<script>
  import API from './api'
  import { waitField } from './effects'

  import Earn from './Earn.svelte'

  import NavEarnIcon from '@/nav_earn.svg'
  import NavShopIcon from '@/nav_shop.svg'
  import NavLeadersIcon from '@/nav_leaders.svg'
  import NavFriendsIcon from '@/nav_friends.svg'
  import NavAdminIcon from '@/tab_active.svg'

  let currentPoints = 0
  let refCount = 0
  let refPoints = 0
  let products = []
  let tasks = []

  let loading = true
  let tutor = null

  let nextProductPoints = 0

  function calcPoints() {
    let p = products.find((p) => p.price > currentPoints)
    nextProductPoints = p ? p.price : currentPoints
  }
  $: calcPoints(currentPoints, products)

  const locs = {
    en: {
      earn: 'Earn',
      shop: 'Shop',
      leads: 'Leaders',
      frens: 'Friends'
    },
    ru: {
      earn: 'Задачи',
      shop: 'Магазин',
      leads: 'Лидеры',
      frens: 'Друзья'
    }
  }
  const loc = locs[LANG] || locs.en

  let pages = [
    { name: loc.earn, icon: NavEarnIcon, counter: 0, component: Earn },
    { name: loc.shop, icon: NavShopIcon, counter: 0 },
    // TODO: { name: loc.leads, icon: NavLeadersIcon, counter: 0 }
    { name: loc.frens, icon: NavFriendsIcon, counter: 0 }
  ]
  let adminPage = { name: 'Admin', icon: NavAdminIcon }
  let currentPage = pages[0]

  function selectPage(page) {
    if (window.Reload) window.Reload()
    currentPage = page
  }

  function updateCounters() {
    pages[0].counter = tasks.filter((t) => t.status == 'claim').length
    pages[1].counter = products.filter((p) => p.price <= currentPoints).length
  }
  $: updateCounters(tasks, products, currentPoints)

  async function loadTutorial() {
    tutor = (await import('./Tutorial.svelte')).default
  }
  function setPage(i, imported) {
    pages[i].component = imported.default
    pages = pages
  }
  async function loadAllModules() {
    setPage(1, await import('./Shop.svelte'))
    setPage(2, await import('./Friends.svelte'))
    // TODO:
    // setPage(3, import('./Leaders.svelte'))
  }
  async function loadAdmin() {
    if (!pages.includes(adminPage)) {
      pages = pages.concat(adminPage)
      setPage(3, await import('./Admin.svelte'))
    }
  }

  fetchOverview()

  if (localStorage.getItem('tutor') != '1') loadTutorial()
  else loadAllModules()

  const tutorDone = () => {
    tutor = null
    loadAllModules()
    API.sendInit().then(() => localStorage.setItem('tutor', '1'))
  }

  function addPoints(delta) {
    currentPoints += delta
    if (currentPoints < 0) currentPoints = 0
  }

  function fetchOverview() {
    API.getOverview(Telegram.WebApp.initDataUnsafe?.start_param || '')
      .then((o) => {
        currentPoints = o.points
        refCount = o.refCount
        refPoints = o.refPoints
        products = o.products
        tasks = o.tasks

        if (o.isAdmin) window.Admin = loadAdmin

        initTasksAdsGram()

        loading = false
        setInterval(checkTasks, 1000)
      })
      .catch(() => setTimeout(fetchOverview, 1000))
  }

  const nowUnix = () => Math.ceil(new Date().getTime() / 1000)
  const startTime = nowUnix()

  let checkNum = 0
  function checkTasks() {
    checkNum++
    if (checkNum % 5 == 0) {
      API.getTasksEvents(startTime).then((events) => {
        if (events?.length) {
          for (let e of events) {
            const t = tasks.find((t) => t.id == e.id && t.subID == e.subID)
            if (t) t.status = e.status
          }
          tasks = tasks
        }
      })
    }

    const nowTime = nowUnix()
    let wasChanges = false
    for (let t of tasks) {
      if (
        (t.status == 'pending' && t.pending && t.start + t.pending < nowTime) ||
        (t.status == 'done' && t.interval && t.start + t.interval < nowTime)
      ) {
        t.status = 'active'
        wasChanges = true
      }
    }
    if (wasChanges) tasks = tasks
  }

  function initTasksAdsGram() {
    if (window.DEBUG) return
    let filtered = tasks.filter((t) =>
      ['ads_gram_rewarded', 'ads_gram_task'].includes(t.type)
    )
    if (filtered.length) {
      filtered.map((t) => (t.loading = true))
      tasks = tasks
      waitField(() => window.Adsgram).then(() => {
        for (let t of filtered) {
          if (t.type == 'ads_gram_rewarded') {
            t.adsGram = Adsgram.init({ blockId: t.actionAdsGramBlockID })
          }
          t.loading = false
        }
        tasks = tasks
      })
    }
  }

  const mobile = navigator.maxTouchPoints > 0
</script>

{#if tutor}
  <svelte:component this={tutor} on:done={tutorDone} />
{:else}
  <div id="app">
    {#if currentPage && currentPage.component}
      <svelte:component
        this={currentPage.component}
        {loading}
        {tasks}
        {products}
        {currentPoints}
        {nextProductPoints}
        {refCount}
        {refPoints}
        {addPoints}
      />
    {/if}
  </div>

  {#if !loading}
    <nav id="nav" class:mobile>
      {#each pages as page}
        <button
          class="material"
          class:active={page == currentPage}
          on:click={() => selectPage(page)}
        >
          <svg>
            <use xlink:href={page.icon + '#i'}></use>
          </svg>
          <span>{page.name}</span>
          <span class="badge" class:show={page.counter}>{page.counter}</span>
        </button>
      {/each}
    </nav>
  {/if}
{/if}
