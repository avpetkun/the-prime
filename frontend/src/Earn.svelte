<script>
  import './earn.scss'

  import API from './api'
  import TON from './ton'
  import Haptic from './haptic'
  import Confetti from './confetti'

  import Progress from './Progress.svelte'
  import Particles from './Particles.svelte'
  import Tabs from './Tabs.svelte'
  import Tasks from './Tasks.svelte'
  import StarIcon from '@/tgstar.svg'

  import getLoc from './loc'
  const loc = getLoc()

  export let loading = true
  export let tasks = []
  export let currentPoints = 0
  export let nextProductPoints = 0
  export let addPoints

  $: tasksClaim = tasks.filter((t) => t.status == 'claim')
  $: tasksSimple = tasks.filter((t) => t.status == 'active' && !t.premium)
  $: tasksPremium = tasks.filter((t) => t.status == 'active' && t.premium)
  $: tasksPending = tasks.filter((t) => t.status == 'pending')
  $: tasksDone = tasks.filter((t) => t.status == 'done')

  const tasksSum = (list) => list.reduce((sum, t) => sum + t.points, 0)

  let tabSelected = 0
  let tabPointsDone = 0
  $: tabPointsActive = tasksSum(tasksClaim)
  $: tabPointsPending = tasksSum(tasksPending)

  let fetching = false
  $: blur = loading || fetching

  const nowUnix = () => Math.ceil(new Date().getTime() / 1000)

  function onTaskStart({ detail }) {
    if (blur) return
    const isActive = detail.status == 'active'
    if (isActive) Haptic.lightInpact()
    executeTask(detail, () => {
      if (isActive) {
        detail.status = 'pending'
        detail.start = nowUnix()
        tasks = tasks

        API.taskStart(detail.id, detail.subID).catch(() => {
          detail.status = 'active'
          detail.start = 0
          tasks = tasks

          Haptic.notifyError()
        })
      }
    })
  }

  function onTaskClaim({ detail }) {
    detail.status =
      detail.interval == 0 || detail.start + detail.interval > nowUnix()
        ? 'done'
        : 'active'
    tasks = tasks

    tabPointsDone += detail.points
    addPoints(detail.points)

    Haptic.notifySuccess()
    Confetti()

    API.taskClaim(detail.id, detail.subID).catch(() => {
      tabPointsDone -= detail.points
      addPoints(-detail.points)
      detail.status = 'claim'
      tasks = tasks

      Haptic.notifyError()
    })
  }

  function openLink(link) {
    if (link.startsWith('https://t.me/') || link.startsWith('tg://'))
      return Telegram.WebApp.openTelegramLink(link)
    Telegram.WebApp.openLink(link)
  }

  function executeTask(t, success) {
    // can click any time
    switch (t.type) {
      case 'free':
        return success()
      case 'invite':
        return execInvite(success)
      case 'join':
      case 'partner_event':
      case 'partner_check':
      case 'free_link':
        success()
        return openLink(t.actionLink)
      case 'ton_disconnect':
        return execTonDisconnect(success)
      case 'tapp_ads':
        success()
        return openLink(t.actionLink)
    }
    if (t.status != 'active') return
    // can click only once
    switch (t.type) {
      case 'ton_connect':
        return execTonConnect(success)
      case 'ton_deposit':
        return execTonDeposit(t.id, success)
      case 'stars_deposit':
        return execStarsDeposit(t.id, success)
      case 'ads_gram_task':
      case 'ads_gram_rewarded':
        return execAdsGram(t, success)
      case 'monetag-link':
      case 'monetag-banner':
        return execMonetag(t, success)
    }
  }

  function execInvite(success) {
    fetching = true
    API.getInviteMessage()
      .then((msgID) => {
        window.Telegram.WebApp.shareMessage(msgID, (sent) => {
          if (sent) success()
        })
      })
      .finally(() => (fetching = false))
  }

  function execTonConnect(success) {
    fetching = true
    TON.getConnect()
      .then(({ tonConnect, tonConnected }) => {
        if (tonConnected) return success()
        tonConnect.connectWallet().then(() => {
          TON.setConnected(true)
          success()
        })
      })
      .finally(() => (fetching = false))
  }

  function execTonDisconnect(success) {
    fetching = true
    TON.getConnect()
      .then(({ tonConnect, tonConnected }) => {
        if (tonConnected)
          tonConnect.disconnect().then(() => {
            TON.setConnected(false)
            success()
          })
      })
      .finally(() => (fetching = false))
  }

  function execTonDeposit(taskID, success) {
    const makeTx = (tonConnect) => {
      fetching = true
      API.getTaskTonInvoice(taskID)
        .then((transaction) => {
          tonConnect.sendTransaction(transaction).then((result) => {
            console.log(result)
            success()
          })
        })
        .finally(() => (fetching = false))
    }
    fetching = true
    TON.getConnect()
      .then(({ tonConnect, tonConnected }) => {
        if (tonConnected) return makeTx(tonConnect)
        tonConnect.connectWallet().then(() => {
          TON.setConnected(true)
          makeTx(tonConnect)
        })
      })
      .finally(() => (fetching = false))
  }

  function execStarsDeposit(taskID, success) {
    fetching = true
    API.getTaskStarsInvoice(taskID)
      .then((invoiceLink) => {
        window.Telegram.WebApp.openInvoice(invoiceLink, (status) => {
          if (status != 'cancelled') success()
        })
      })
      .finally(() => (fetching = false))
  }

  //

  function execAdsGram(t, success) {
    if (!t.adsGram) return success()
    t.adsGram.show().then((result) => {
      if (result.done) success()
    })
  }

  function execMonetag(t, success) {
    const opts =
      t.type == 'monetag-link'
        ? { ymid: `${TG_USER}-${t.id}`, type: 'pop' }
        : { ymid: `${TG_USER}-${t.id}` }
    show_9031733(opts).then(success)
  }
</script>

<div class="earn pop-page">
  <div class="header fetcher" class:blur>
    <Particles />
    <Progress current={currentPoints} total={nextProductPoints} />
    <div class="text">
      <span class="title">{loc.top.title}</span>
      <div class="desc">
        <span>{@html loc.top.desc1}</span>
        <span>{@html loc.top.desc2}</span>
      </div>
    </div>
  </div>

  {#if loading || fetching}
    <img src={StarIcon} alt="" />
  {/if}

  <div class="fetcher" class:blur>
    <Tabs
      bind:value={tabSelected}
      countActive={tabPointsActive}
      countPending={tabPointsPending}
      countDone={tabPointsDone}
    />

    {#if !loading}
      {#if tasksClaim.length}
        <div class="tasks-block">
          <div class="block-title">{loc.blocks.claim}</div>
          <Tasks claim tasks={tasksClaim} on:claim={onTaskClaim} />
        </div>
      {/if}

      {#if tabSelected == 0}
        <div class="pop-page">
          {#if tasksPremium.length}
            <div class="tasks-block">
              <div class="block-title">{loc.blocks.premium}</div>
              <Tasks premium tasks={tasksPremium} on:start={onTaskStart} />
            </div>
          {/if}
          <div class="tasks-block">
            <div class="block-title">{loc.blocks.simple}</div>
            <Tasks tasks={tasksSimple} on:start={onTaskStart} />
          </div>
        </div>
      {:else if tabSelected == 1}
        <div class="pop-page">
          <div class="tasks-block">
            <div class="block-title">{loc.blocks.pending}</div>
            <Tasks tasks={tasksPending} on:start={onTaskStart} />
          </div>
        </div>
      {:else}
        <div class="pop-page">
          <div class="tasks-block">
            <div class="block-title">{loc.blocks.done}</div>
            <Tasks tasks={tasksDone} on:start={onTaskStart} />
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>
