<script>
  import './progress.scss'
  import { startFollow } from './effects'
  import Haptic from './haptic'

  export let current = 0
  export let total = 0

  $: value = total ? Math.min(1, current / total) : 0
  $: track = (value + 1) / 2

  // 420 = 0%, 154 = 100%
  const circle = (v) => 154 + Math.floor((1 - v) * 266)
  $: valueCircle = circle(value)
  $: trackCircle = circle(track)

  let currentAnim = 0
  let totalAnim = 0

  function animateValues() {
    const [current0, total0] = [currentAnim, totalAnim]
    startFollow(500, (follow) => {
      currentAnim = follow(current0, current)
      totalAnim = follow(total0, total)
    })
  }
  $: animateValues(current, total)

  let rotate = false
  let rollCount = 0
  function doRoll() {
    rotate = !rotate
    Haptic.lightInpact()
    rollCount++
    if (rollCount == 5 && window.Admin) window.Admin()
  }
</script>

<div class="progress" class:rotate on:pointerdown={doRoll}>
  <div class="circle">
    <svg>
      <defs>
        <linearGradient id="Grad">
          <stop offset="0%" stop-color="#3E8BFF"></stop>
          <stop offset="100%" stop-color="#847CFF"></stop>
        </linearGradient>
      </defs>
      <circle
        stroke="url(#Grad)"
        stroke-dashoffset={trackCircle}
        opacity="0.15"
        cx="52.5"
        cy="52.5"
        r="42.5"
      />
      <circle
        stroke="url(#Grad)"
        stroke-dashoffset={valueCircle}
        cx="52.5"
        cy="52.5"
        r="42.5"
      />
    </svg>
    <div class="legend">
      <span class:zoom={current != currentAnim}>{currentAnim}</span>
      <i></i>
      <span class:zoom={total != totalAnim}>{totalAnim}</span>
    </div>
  </div>
</div>
