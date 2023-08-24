<script lang="ts">
  export let name = 'name'
  export let endpointInput = ''
  let loading = false
  let pingFeedback = ''
  let pingSuccess = -1

  const pingServer = async () => {
    if (endpointInput.length == 0 || endpointInput === null || endpointInput === undefined) return;
    loading = true;
    try {
       await fetch(endpointInput)
      pingFeedback = 'input-solid-success'
      pingSuccess = 1
    } catch(err) {
      pingFeedback = 'input-solid-error'
      pingSuccess = 0
    } finally {
      loading = false
    }
  }
</script>
<div class="w-full">
  <label class="form-label" for="name">{name}</label>
  <div class="flex flex-row">
    <input bind:value={endpointInput} class="{pingFeedback} input input-solid max-w-full" placeholder="E.g http://localhost:8000" type="text" id="name"/>
    <span class="tooltip tooltip-left" data-tooltip="Check connection">
      {#if loading}
        <svg class="spinner-ring w-[70px] [--spinner-color:var(--slate-12)]" viewBox="25 25 50 50" stroke-width="5">
          <circle cx="50" cy="50" r="10"/>
        </svg>
      {:else}
        <button on:click={() => pingServer()} class=" w-[70px] btn btn-primary">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
            <path d="M3.478 2.405a.75.75 0 00-.926.94l2.432 7.905H13.5a.75.75 0 010 1.5H4.984l-2.432 7.905a.75.75 0 00.926.94 60.519 60.519 0 0018.445-8.986.75.75 0 000-1.218A60.517 60.517 0 003.478 2.405z"/>
          </svg>
        </button>
      {/if}
    </span>
  </div>
  {#if pingSuccess == 1}
    <div class="form-label">
      <span class="form-label-alt text-success">Success</span>
    </div>
    {:else if pingSuccess == 0}
    <div class="form-label">
      <span class="form-label-alt text-error">Error</span>
    </div>
  {/if}
</div>
