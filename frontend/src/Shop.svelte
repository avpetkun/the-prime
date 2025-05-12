<script>
  import './shop.scss'

  import API from './api'
  import Haptic from './haptic'
  import Confetti from './confetti'
  import { startFollow } from './effects'

  import Particles from './Particles.svelte'

  import TgStarIcon from '@/tgstar.svg'
  import TgPrimeIcon from '@/prime_96.png'

  function typeIcon(type) {
    return {
      tg_stars: TgStarIcon,
      tg_premium: TgPrimeIcon
    }[type]
  }

  const locs = {
    en: {
      buy: 'Buy',
      noUsername: 'Please, set username for account',
      notEnough: 'Insufficient points',
      claimOk: 'Congratulations!',
      claimErr: 'Error. Try later',
      yourPoints: 'Your points',
      spendInfo: 'Spend your points on rewards!',
      remaining0: 'Earn',
      remaining1: 'points to claim'
    },
    ru: {
      buy: 'Купить',
      noUsername: 'Пожалуйста, задайте юзернейм для аккаунта',
      notEnough: 'Недостаточно очков',
      claimOk: 'Поздравляем!',
      claimErr: 'Ошибка. Попробуйте позже',
      yourPoints: 'Твой баланс',
      spendInfo: 'Трать свои баллы на призы!',
      remaining0: 'Осталось',
      remaining1: 'очков'
    }
  }
  const loc = locs[LANG] || locs.en

  export let products = []
  export let currentPoints = 0
  export let addPoints

  let currentAnim = 0

  function animateValues() {
    const [current0] = [currentAnim]
    startFollow(500, (follow) => {
      currentAnim = follow(current0, currentPoints)
    })
  }
  $: animateValues(currentPoints)

  let fetching = false

  function fmtNumber(n) {
    return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  }
  $: fmtPoints = fmtNumber(currentAnim)

  function enough(p) {
    return p.price <= currentPoints
  }

  function remaining(p) {
    let r = p.price - currentPoints
    if (r > 9999) r = `${Math.floor(r / 1000)}k`
    return `${loc.remaining0} ${r} ${loc.remaining1}`
  }

  function claim(p) {
    if (fetching || p.price <= 0 || p.amount <= 0) return

    Haptic.lightInpact()
    if (!enough(p)) {
      // console.warn(loc.notEnough)
      Telegram.WebApp.showAlert(loc.notEnough)
      Haptic.notifyError()
      return
    }
    function exec() {
      if (!window.DEBUG && ['tg_stars', 'tg_premium'].includes(p.type)) {
        if (!Telegram.WebApp.initDataUnsafe.user.username) {
          Telegram.WebApp.showAlert(loc.noUsername)
          Haptic.notifyError()
          return
        }
      }
      fetching = true
      API.productClaim(p.id)
        .then((spendPoints) => {
          if (p.price != spendPoints) p.price = spendPoints
          addPoints(-spendPoints)
          console.info(loc.claimOk)
          Telegram.WebApp.showAlert(loc.claimOk)
          Haptic.notifySuccess()
          Confetti()
        })
        .catch(() => {
          console.error(loc.claimErr)
          Telegram.WebApp.showAlert(loc.claimErr)
          Haptic.notifyError()
        })
        .finally(() => (fetching = false))
    }
    if (window.DEBUG) return confirm(`${loc.buy} ${p.name}?`) && exec()
    Telegram.WebApp.showConfirm(`${loc.buy} ${p.name}?`, (ok) => ok && exec())
  }
</script>

<div class="shop pop-page">
  <Particles horizontal vertical count={300} />

  <div class="header">
    <div class="texts">
      <span class="title">{loc.yourPoints}</span>
      <span class="points">{fmtPoints}</span>
      <span>{loc.spendInfo}</span>
    </div>
  </div>

  {#if fetching}
    <img src={TgStarIcon} alt="" />
  {/if}

  <div class="products fetcher" class:fetching>
    {#each products as p}
      <button class="product material" on:click={() => claim(p)}>
        <div class="icon">
          <img src={typeIcon(p.type)} alt="" />
        </div>
        <div class="content">
          <div class="name">
            <span>{p.name}</span>
            {#if p.badge}
              <span class="badge">{p.badge}</span>
            {/if}
          </div>
          {#if !enough(p)}
            <span class="info">{remaining(p, currentPoints)}</span>
          {/if}
        </div>
        {#if p.price}
          <div class="points" class:active={enough(p)}>
            {fmtNumber(p.price)}
          </div>
        {/if}
      </button>
    {/each}
  </div>
</div>
