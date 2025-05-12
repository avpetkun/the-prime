<script>
  import './friends.scss'

  import QRCode from 'qrcode'
  import API from './api'
  import Haptic from './haptic'
  import { startFollow } from './effects'

  import Particles from './Particles.svelte'
  import Progress from './Progress.svelte'
  import StarIcon from '@/tgstar.svg'
  import ArrowIcon from '@/friends_arrow.svg'
  import QrCodeIcon from '@/qr.svg'

  const locs = {
    en: {
      title: 'Referrals',
      desc: 'The amount of people you invited',
      numFriends: 'friends',
      numPoints: 'points',
      invite1: 'Invite more friends',
      invite2: 'Click here to share your referral link'
    },
    ru: {
      title: 'Рефералы',
      desc: 'Люди которых ты пригласил',
      numFriends: 'людей',
      numPoints: 'баллов',
      invite1: 'Пригласи больше друзей',
      invite2: 'Кликни чтобы поделиться ссылкой'
    }
  }
  const loc = locs[LANG] || locs.en

  export let currentPoints = 0
  export let nextProductPoints = 0
  export let refCount = 0
  export let refPoints = 0

  let refPointsAnim = 0
  let refCountAnim = 0

  function animateValues() {
    const [refp0, refc0] = [refPointsAnim, refCountAnim]
    startFollow(500, (follow) => {
      refPointsAnim = follow(refp0, refPoints)
      refCountAnim = follow(refc0, refCount)
    })
  }
  $: animateValues(refPoints, refCount)

  function fmtNumber(n) {
    return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  }
  $: fmtRefPoints = fmtNumber(refPointsAnim)
  $: fmtRefCount = fmtNumber(refCountAnim)

  let rotate = false
  let qrURL = ''

  function onQR() {
    if (fetching) return
    rotate = !rotate
    if (!qrURL) {
      QRCode.toDataURL(
        `https://t.me/YourPrimeBot/get?startapp=${TG_USER}`,
        (err, url) => {
          if (err) throw err
          qrURL = url
        }
      )
    }
  }

  let fetching = false
  let shareMessageID = 0

  function onShare() {
    if (fetching) return

    function share() {
      window.Telegram.WebApp.shareMessage(shareMessageID)
    }
    Haptic.lightInpact()
    if (shareMessageID) return share()

    fetching = true
    API.getInviteMessage()
      .then((msgID) => {
        shareMessageID = msgID
        share()
      })
      .finally(() => (fetching = false))
  }
</script>

<div class="friends pop-page">
  {#if fetching}
    <img src={StarIcon} alt="" />
  {/if}

  <div class="header">
    <Particles horizontal vertical count={300} />
    <Progress current={currentPoints} total={nextProductPoints} />
    <span class="title">{loc.title}</span>
    <span class="desc">{loc.desc}</span>
  </div>

  <div class="card fetcher" class:fetching>
    <div class="inner" class:rotate>
      <div class="front">
        <div class="text top">
          <span>{fmtRefCount}</span>
          <span>{loc.numFriends}</span>
        </div>
        <div class="text bottom">
          <span>{fmtRefPoints}</span>
          <span>{loc.numPoints}</span>
        </div>
        <div class="arrow">
          <svg>
            <use xlink:href={ArrowIcon + '#i'}></use>
          </svg>
        </div>
      </div>
      <div class="back">
        <button on:click={onShare}>
          <img src={qrURL} alt="" />
        </button>
      </div>
    </div>
  </div>

  <button class="material invite fetcher" class:fetching on:click={onShare}>
    <div class="text">
      <span>{loc.invite1}</span>
      <span>{loc.invite2}</span>
    </div>
    <button class="material qr" on:click|stopPropagation={onQR}>
      <img src={QrCodeIcon} alt="" />
    </button>
  </button>
</div>
