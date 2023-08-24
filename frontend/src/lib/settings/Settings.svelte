<script>
  import {slide} from 'svelte/transition'
  import {onMount} from 'svelte'
  import {scale} from 'svelte/transition'
  import Modal from "../components/Modal.svelte";
  import ServerEndpointInput from "./ServerEndpointInput.svelte";
  import {showSettingsModal as showModal} from '../stores/settings'

  let activeTab = 0
  const activeTabClasses = ' tab-underline tab-active '

  let mainServerEndpoint = ''
  let preprocessingServerEndpoint = ''
  let storageServerEndpoint = ''

  onMount(() => {
    mainServerEndpoint = localStorage.getItem("settings:endpoints:main") || ''
    preprocessingServerEndpoint = localStorage.getItem("settings:endpoints:preprocessing") || ''
    storageServerEndpoint = localStorage.getItem("settings:endpoints:storage") || ''
  })

  const saveSettings = () => {
    localStorage.setItem("settings:endpoints:main", mainServerEndpoint);
    localStorage.setItem("settings:endpoints:preprocessing", preprocessingServerEndpoint);
    localStorage.setItem("settings:endpoints:storage", storageServerEndpoint);

    showModal.set(false)
  }
</script>

<Modal show={showModal} title="Settings">
  <div class="settings-form overflow-hidden">
    <div class="tabs w-full">
      <button on:click={() => activeTab = 0}
        class="{activeTab === 0 || activeTab > 2 ? activeTabClasses : ''} tab px-6">
        Server Endpoints
      </button>
      <button on:click={() => activeTab = 1} class="{activeTab === 1 ? activeTabClasses : '' } tab px-6">Tab 2</button>
      <button on:click={() => activeTab = 2} class="{activeTab === 2 ? activeTabClasses : '' } tab px-6">Tab 3</button>
    </div>
    <div class="tabs-content my-4 lg:p-3">
      {#if activeTab === 1}
        <div in:scale={{start: 3, opacity: 0.5}}>
          <h1>Tab2</h1>
        </div>
        {:else if activeTab === 2}
        <div in:scale={{start: 3, opacity: 0.5}}>
          <h1>Tab3</h1>
        </div>
      {:else}
        <div class="flex flex-col" style="gap: 10px" in:scale={{start: 3, opacity: 0.5}}>
          <ServerEndpointInput bind:endpointInput={mainServerEndpoint} name="Main Service"/>
          <ServerEndpointInput bind:endpointInput={preprocessingServerEndpoint} name="Preprocessing Service"/>
          <ServerEndpointInput bind:endpointInput={storageServerEndpoint} name="Storage Service"/>
        </div>
      {/if}
    </div>
  </div>
  <div class="flex gap-3">
    <button class="btn btn-success btn-block" on:click={() => saveSettings()}>Save</button>
    <button class="btn btn-block" on:click={() => $showModal = false}>Close</button>
  </div>
</Modal>
