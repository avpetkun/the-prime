<script>
  import './task.scss'

  import { createEventDispatcher } from 'svelte'
  import { waitField } from './effects'

  import PrimeIcon from '@/prime_26.png'
  import RefreshIcon from '@/refresh.svg'

  import getLoc from './loc'
  const loc = getLoc()

  export let premium = false

  export let type = ''
  export let icon = ''
  export let name = ''
  export let desc = ''
  export let status = ''
  export let points = 0
  export let claim = false
  export let loading = false

  let show = true

  let i = ['monetag-link', 'monetag-banner'].indexOf(type)
  if (i != -1 && status == 'active') {
    loading = true
    waitField(() => window.show_9031733).then(() => {
      if (i == 0) return (loading = false)
      show_9031733({ type: 'preload', timeout: 10 })
        .then(() => (loading = false))
        .catch(() => (show = false))
    })
  }

  const dispatch = createEventDispatcher()
  const onStart = () => !loading && dispatch('start')
  const onClaim = () => dispatch('claim')
</script>

{#if show}
  <button class="task material" class:premium on:click={onStart}>
    <div class="icon">
      <img src={icon || PrimeIcon} alt="" />
    </div>
    <div class="text">
      <div class="name">{name}</div>
      {#if desc}
        <div class="desc">{desc}</div>
      {/if}
    </div>
    <div class="points" class:claim>
      {#if loading}
        <img src={RefreshIcon} alt="" />
      {:else if claim}
        <button class="material dark" on:click|stopPropagation={onClaim}>
          {loc.claim}
        </button>
      {:else if points}
        <span>+ {points}</span>
      {/if}
    </div>
  </button>
{/if}
