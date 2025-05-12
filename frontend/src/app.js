import './app.scss'
import './material'

{
  const w = Telegram.WebApp
  window.LANG ||= w.initDataUnsafe.user?.language_code
  window.AUTH ||= w.initData
  window.TG_USER ||= w.initDataUnsafe.user.id

  function init() {
    document.documentElement.className = window.SCHEME || w.colorScheme

    const t = window.THEME || w.themeParams
    let bg = t.bg_color
    if (bg == t.section_bg_color) {
      bg = t.secondary_bg_color
    }
    w.backgroundColor = bg
    w.headerColor = bg

    document.documentElement.style.background = bg
    document.documentElement.style.setProperty('--tg-theme-bg-color', bg)
  }
  w.onEvent('themeChanged', init)
  init()
  w.disableVerticalSwipes()
  w.expand()
}

import App from './App.svelte'

const app = new App({
  target: document.body
})

export default app
