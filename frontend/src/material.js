import './material.scss'

let lights = []

function start(target, cursor) {
  if (!target.classList.contains('material')) return

  const r = target.getBoundingClientRect()
  const ox = cursor.clientX - r.left
  const oy = cursor.clientY - r.top

  const el = document.createElement('i')
  el.className = 'material-spot'
  el.style.left = ox + 'px'
  el.style.top = oy + 'px'

  const size = Math.max(r.width, r.height) * 2
  el.style.width = size + 'px'
  el.style.height = size + 'px'

  lights.push(el)
  target.appendChild(el)
}

function stop() {
  const oldLights = lights
  lights = []
  for (let el of oldLights) el.classList.add('material-fade')
  setTimeout(() => oldLights.map((el) => el.remove()), 100)
}

if (navigator.maxTouchPoints > 0) {
  document.addEventListener('touchstart', (e) => start(e.target, e.touches[0]))
  document.addEventListener('touchend', stop)
} else {
  document.addEventListener('mousedown', (e) => start(e.target, e))
  document.addEventListener('mouseup', stop)
}
