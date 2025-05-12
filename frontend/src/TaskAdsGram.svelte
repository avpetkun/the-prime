<script>
  import './TaskAdsGram.scss'
  import { onDestroy, createEventDispatcher } from 'svelte'

  import getLoc from './loc'
  const loc = getLoc()

  export let premium = false
  export let points = 0
  export let loading = false
  export let blockID = 0

  const dispatch = createEventDispatcher()
  const onStart = () => !loading && dispatch('start')

  let task = null
  let show = true

  function onCancel() {
    show = false
  }

  $: if (task) {
    task.addEventListener('reward', onStart)
    task.addEventListener('onBannerNotFound', onCancel)
  }
  onDestroy(() => {
    if (task) {
      task.removeEventListener('reward', onStart)
      task.removeEventListener('onBannerNotFound', onCancel)
    }
  })
</script>

{#if show}
  <adsgram-task
    data-block-id={blockID}
    class="task task-adsgram"
    class:premium
    bind:this={task}
  >
    <div slot="reward" class="reward">+ {points} {loc.points}</div>
    <div slot="button" class="button material dark">Go</div>
  </adsgram-task>
{/if}
