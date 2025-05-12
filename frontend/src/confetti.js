import './confetti.css'
import TgStarIcon from '@/tgstar.svg'

function launch() {
  let elems = []
  let nline = Math.floor(screen.availWidth / 30 + 1)
  for (let i = 0; i < 80; i++) {
    let el = document.createElement('img')
    el.src = TgStarIcon

    let x = (i % nline) * 30
    let size = Math.min(1, Math.random() * 2) * 30
    let delay = Math.random() * 0.66
    let start = i % 2 == 0 ? '-start' : ''

    el.className = `confetti`
    el.style = `animation-name: blink${start}, confetti; animation-delay: ${delay}s; left:${x}px; width:${size}px; height:${size}px`

    document.body.appendChild(el)
    elems.push(el)
  }
  setTimeout(() => elems.map((el) => document.body.removeChild(el)), 2000)
}

export default launch
