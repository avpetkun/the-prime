<script>
  import './particles.scss'
  import ParticleA from '@/particleA.svg'
  import ParticleB from '@/particleB.svg'

  export let horizontal = false
  export let vertical = false
  export let count = 150

  let particles = []

  function emit(n) {
    for (let i = 0; i < n; i++) {
      const y = vertical
        ? (Telegram.WebApp.viewportHeight - 53) * Math.random()
        : horizontal
          ? 138 * Math.random() // height 150px - 12px size
          : 12 + 90 * Math.random() // 64 center, 12 top margin, 105 progress height
      const size = Math.min(1, Math.random() * 2)
      let side = i < n * 0.5 ? 'left' : 'right'
      let vert = (i < n * 0.5 ? i < n * 0.25 : i < n * 0.75) ? '-up' : '-down'
      let start = i % 2 == 0 ? '-start' : ''
      if (horizontal) {
        side = 'left'
        vert = start = ''
      }
      const anim = `animation-name: blink${start}, move-${side}${start}${vert};`
      const delay = Math.random() * 4
      particles = particles.concat({
        start,
        style: `${anim} top:${y}px; width:${12 * size}px; height:${16 * size}px; animation-delay:${delay}s`
      })
    }
  }
  emit(count)
</script>

<div class="particles" class:horizontal>
  {#each particles as p}
    <img src={p.start ? ParticleA : ParticleB} alt="" style={p.style} />
  {/each}
</div>
