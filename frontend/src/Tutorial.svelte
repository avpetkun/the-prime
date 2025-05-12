<script>
  import './tutorial.scss'

  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  import BgStars from '@/bg_stars.svg'
  import TgStar from '@/tgstar.svg'

  import TutorIcon1 from '@/tutor-1.svg'
  import TutorIcon2 from '@/tutor-2.svg'
  import TutorIcon3 from '@/tutor-3.svg'
  import TutorArrow from '@/tutor-arrow.svg'

  const locs = {
    en: {
      title: 'Welcome to Premium Space!',
      item1: 'Complete tasks – get points!',
      item2: 'Exchange points for Telegram Premium and Stars Packs!',
      item3: 'Wait for updates, and exchange points for other things!',
      okay: 'Okay'
    },
    ru: {
      title: 'Добро пожаловать в Premium Space!',
      item1: 'Выполняй задачи – зарабатывай баллы!',
      item2: 'Обменивай баллы на Telegram Premium и Stars!',
      item3: 'Ждите обновлений и обменивайте очки на другие награды!',
      okay: 'Хорошо'
    }
  }
  const loc = locs[LANG] || locs.en

  const button = Telegram.WebApp.MainButton

  function skip() {
    dispatch('done')
    button.hide()
  }

  button.onClick(skip)
  button.setParams({
    text: loc.okay,
    has_shine_effect: true,
    is_active: true,
    is_visible: true,
    color: '#5486fe',
    text_color: '#fff'
  })
</script>

<div class="tutor">
  <img src={BgStars} alt="" class="bg" />

  <div class="star">
    <img src={TgStar} alt="" />
  </div>

  <div class="title">{loc.title}</div>

  <div class="item">
    <img src={TutorIcon1} alt="" />
    <span>{loc.item1}</span>
  </div>

  <img src={TutorArrow} alt="" class="arrow" />

  <div class="item">
    <img src={TutorIcon2} alt="" />
    <span>{loc.item2}</span>
  </div>

  <img src={TutorArrow} alt="" class="arrow" />

  <div class="item">
    <img src={TutorIcon3} alt="" />
    <span>{loc.item3}</span>
  </div>
</div>

{#if window.DEBUG}
  <button on:click={skip}>Skip tutorial</button>
{/if}
