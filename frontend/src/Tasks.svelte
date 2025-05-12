<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  import Task from './Task.svelte'
  import TaskAdsGram from './TaskAdsGram.svelte'

  import getLoc from './loc'
  const loc = getLoc()

  export let premium = false
  export let claim = false
  export let tasks = []

  function isAdsGramActiveTask(t) {
    return t.type == 'ads_gram_task' && t.status == 'active' && !t.loading
  }
</script>

{#each tasks as t (t.id + ':' + t.subID)}
  {#if isAdsGramActiveTask(t)}
    <TaskAdsGram
      {premium}
      points={t.points}
      loading={t.loading}
      blockID={t.actionAdsGramBlockID}
      on:start={() => dispatch('start', t)}
    />
  {:else}
    <Task
      {premium}
      {claim}
      type={t.type}
      icon={t.icon}
      name={t.name}
      desc={t.desc}
      status={t.status}
      points={t.points}
      loading={t.loading}
      on:start={() => dispatch('start', t)}
      on:claim={() => dispatch('claim', t)}
    />
  {/if}
{:else}
  <div class="tasks-empty">{loc.noTasks}</div>
{/each}
