<script lang="ts">
  import { writable } from 'svelte/store';
  import {imagesArrayData} from '../../lib/stores/images'
  import ImageDetailsModal from '../components/Modal.svelte';
  import type {Image} from '../../../wailsjs/go/models'
  import { OpenImage } from "../../../wailsjs/go/main/App"


  let imagesActivePage = 0
  let pages = []

  const paginateArray = (pageSize = 10) => {
    var newArr = []
    var maxLength = $imagesArrayData.length
    var newStart = 0
    for (let i = 0; i < Math.round(maxLength/ 10); i++) {
      newStart = i * pageSize
      newArr.push($imagesArrayData.slice(newStart, newStart + pageSize))
    }
    return newArr
  }

  let pageButtonEditable = false
  $: pageInputValue = imagesActivePage

  $: {
    if (pageInputValue >= pages.length ) pageInputValue = pages.length
    if (pageInputValue <= 1) pageInputValue = 2
  }

  const pageInputHandler = (e: KeyboardEvent) => {
    console.log(e.key)
    switch(e.key) {
      case 'Escape':
        pageInputValue = imagesActivePage + 1
        pageButtonEditable = false
        break;
      case 'Enter':
        imagesActivePage = pageInputValue - 1
        pageButtonEditable = false
    }
  }


  imagesArrayData.subscribe(() => {
    pages = paginateArray()
  })

  let imageDetailsModalState = writable(false)
  let imageDetailsSelected: Image
  const selectImage = (img: Image) => {
      imageDetailsSelected = img
      imageDetailsModalState.set(true)
    }
</script>
{#if $imagesArrayData.length <= 0}
  <div class="flex flex-col items-center justify-center w-full min-h-[500px] opacity-50">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-40 h-40">
      <path fill-rule="evenodd" d="M10.5 3.75a6.75 6.75 0 100 13.5 6.75 6.75 0 000-13.5zM2.25 10.5a8.25 8.25 0 1114.59 5.28l4.69 4.69a.75.75 0 11-1.06 1.06l-4.69-4.69A8.25 8.25 0 012.25 10.5z" clip-rule="evenodd" />
    </svg>
    <h1 class="text-2xl">
      Not data loaded yet
    </h1>
  </div>
{:else}
  <!-- Table -->
  <div class="px-4 py-2">
    <h1>Total Items: {$imagesArrayData.length}</h1>
    <div class="flex w-full overflow-x-auto">
      <table class="table-hover table">
        <thead>
          <tr>
            <th>Photo ID</th>
            <th>Title</th>
            <th>Description</th>
            <th>Photo URL</th>
            <th>Tags</th>
            <th>Categories</th>
          </tr>
        </thead>
        <tbody>
          {#each pages[imagesActivePage] as i}
            <tr class="break-words" on:click={() => selectImage(i)}>
              <td>{i.id}</td>
              <td>{i.title}</td>
              <td>{i.description}</td>
              <td>
                <button class="link link-underline" on:click={() => OpenImage(i.image_url)}>
                  {i.image_url}
                </button>
              </td>
              <td>{i.tags.join(',')}</td>
              <td>{i.category[0]}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    <div class="pagination mt-4">
      <span class="tooltip tooltip-top" data-tooltip="Previous Page">
        <button disabled={imagesActivePage === 0} on:click={() => imagesActivePage -= 1} class="btn">
          <svg width="18" height="18" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M12.2574 5.59165C11.9324 5.26665 11.4074 5.26665 11.0824 5.59165L7.25742 9.41665C6.93242 9.74165 6.93242 10.2667 7.25742 10.5917L11.0824 14.4167C11.4074 14.7417 11.9324 14.7417 12.2574 14.4167C12.5824 14.0917 12.5824 13.5667 12.2574 13.2417L9.02409 9.99998L12.2574 6.76665C12.5824 6.44165 12.5741 5.90832 12.2574 5.59165Z" fill="#969696" />
          </svg>
        </button>
      </span>
      <span class="tooltip tooltip-top" data-tooltip="First Page">
        <button class="btn {imagesActivePage === 0 ? 'btn-active' : ''}" on:click={() => imagesActivePage = 0} >{1}</button>
      </span>
      {#if pageButtonEditable}
        <input on:keydown={pageInputHandler} max={pages.length} min="0" type="number" class="input w-24" bind:value={pageInputValue} />
      {:else}
        <button on:click={() => {pageButtonEditable = true}} class="btn { imagesActivePage > 0 && imagesActivePage < pages.length - 1 ? 'btn-active' : ''}" on:click={() => imagesActivePage = 0} >{imagesActivePage > 0 && imagesActivePage < pages.length - 1 ? imagesActivePage + 1: '...'}</button>
      {/if}
      <span class="tooltip tooltip-top" data-tooltip="Last Page">
        <button class="btn {imagesActivePage === pages.length - 1 ? 'btn-active' : ''}" on:click={() => imagesActivePage = pages.length - 1} >{pages.length}</button>
      </span>
      <span class="tooltip tooltip-top" data-tooltip="Next Page">
        <button disabled={imagesActivePage + 1 === pages.length} on:click={() => imagesActivePage += 1} class="btn">
          <svg width="18" height="18" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" clip-rule="evenodd" d="M7.74375 5.2448C7.41875 5.5698 7.41875 6.0948 7.74375 6.4198L10.9771 9.65314L7.74375 12.8865C7.41875 13.2115 7.41875 13.7365 7.74375 14.0615C8.06875 14.3865 8.59375 14.3865 8.91875 14.0615L12.7437 10.2365C13.0687 9.91147 13.0687 9.38647 12.7437 9.06147L8.91875 5.23647C8.60208 4.9198 8.06875 4.9198 7.74375 5.2448Z" fill="#969696" />
          </svg>
        </button>
      </span>
    </div>
  </div>

  <!-- {#if $imageDetailsModalState} -->
  <!--   <ImageDetailsModal title={imageDetailsSelected.title} show={imageDetailsModalState}> -->
  <!--     Body -->
  <!--   </ImageDetailsModal> -->
  <!-- {/if} -->
{/if}

