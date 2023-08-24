<script lang="ts">
  import {imagesArrayData} from './lib/stores/images'
  import Header from './lib/components/Header.svelte'
  import Settings from "./lib/settings/Settings.svelte";
  import { loadingFileIndicator } from './lib/stores/indicators'
  import { onMount } from 'svelte';
  import { cacheLoading  } from './lib/stores/indicators'
  import type {Image} from '../wailsjs/go/models'
  import autoAnimate from '@formkit/auto-animate'

  let cacheStatusText = ""
  onMount(() => {
    window.runtime.EventsOn("cache-start", () => {
      cacheLoading.set(true)
    })

    window.runtime.EventsOn("cache-end", () => {
      cacheLoading.set(false)
    })
    var index = 0
    window.runtime.EventsOn("cache-event", (d: {image: Image, totalImages: number}) => {
      cacheStatusText = `Caching Image ${d.image.id}. (${index}/${d.totalImages})`
      index++
    })
  })
</script>

{#if $cacheLoading}
  <div class="w-full h-screen" use:autoAnimate>
    <div class="flex flex-col p-10 text-center w-full h-full items-center justify-center">
      <h1 class="text-2xl">
        Caching Loaded data. Please wait. 
        This can take minutes to hours based on given data
      </h1>
      <p class="my-5">
        {cacheStatusText}
      </p>
    </div>
  </div>
{:else}
  <Header />
  <div class="">
    <div class="content">
      {#if $loadingFileIndicator}
        <div class="w-full min-h-[400px] flex flex-col items-center justify-center">
          <div class="spinner-wave">
            <div class="spinner-wave-dot w-4 h-4"></div>
            <div class="spinner-wave-dot w-4 h-4"></div>
            <div class="spinner-wave-dot w-4 h-4"></div>
            <div class="spinner-wave-dot w-4 h-4"></div>
          </div>
          <h1 class="text-2xl opacity-50 my-5">Loading file...</h1>
        </div>
      {:else}
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

            <div class="flex w-full overflow-x-auto">
              <table class="table-zebra table">
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
                  {#each $imagesArrayData as i}
                    <tr class="break-words">
                      <td>{i.id}</td>
                      <td>{i.title}</td>
                      <td>{i.description}</td>
                      <td>{i.image_url}</td>
                      <td>{i.tags.join(',')}</td>
                      <td>{i.category.join(',')}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            <div class="pagination mt-4">
              <button class="btn">
                <svg width="18" height="18" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path fill-rule="evenodd" clip-rule="evenodd" d="M12.2574 5.59165C11.9324 5.26665 11.4074 5.26665 11.0824 5.59165L7.25742 9.41665C6.93242 9.74165 6.93242 10.2667 7.25742 10.5917L11.0824 14.4167C11.4074 14.7417 11.9324 14.7417 12.2574 14.4167C12.5824 14.0917 12.5824 13.5667 12.2574 13.2417L9.02409 9.99998L12.2574 6.76665C12.5824 6.44165 12.5741 5.90832 12.2574 5.59165Z" fill="#969696" />
                </svg>
              </button>
              {#each {length: 3} as _, i }
                <button class="btn">1</button>
              {/each}
              <!-- <button class="btn {imagesPageActive === (imagesPageActive + 1)? 'btn-active' : ''}" on:click={() => setNewPage(imagesPageActive + 1)}>{determineNextPage(imagesPageActive , 2)}</button> -->
              <!-- <button class="btn {imagesPageActive === (imagesPageActive + 2)? 'btn-active' : ''}" on:click={() => setNewPage(imagesPageActive + 2)}>{determineNextPage(imagesPageActive , 3)}</button> -->
              <button class="btn">
                <svg width="18" height="18" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path fill-rule="evenodd" clip-rule="evenodd" d="M7.74375 5.2448C7.41875 5.5698 7.41875 6.0948 7.74375 6.4198L10.9771 9.65314L7.74375 12.8865C7.41875 13.2115 7.41875 13.7365 7.74375 14.0615C8.06875 14.3865 8.59375 14.3865 8.91875 14.0615L12.7437 10.2365C13.0687 9.91147 13.0687 9.38647 12.7437 9.06147L8.91875 5.23647C8.60208 4.9198 8.06875 4.9198 7.74375 5.2448Z" fill="#969696" />
                </svg>
              </button>
            </div>

          </div>
        {/if}
      {/if}
      <div>
        <!-- logs  -->
      </div>

      <!---->
      <!-- <div class="flex flex-col" style="gap: 20px;"> -->
      <!--   <ResizableBox direction="bottom" class="w-full bg-green-700 rounded"> -->
      <!--     <h1>ResizableBox</h1> -->
      <!--     <p>Bottom</p> -->
      <!--   </ResizableBox> -->
      <!--   <ResizableBox direction="top" class="w-full bg-green-700"> -->
      <!--     <h1>ResizableBox</h1> -->
      <!--     <p>Top</p> -->
      <!--   </ResizableBox> -->
      <!--   <ResizableBox direction="left" class="w-full bg-green-700"> -->
      <!--     <h1>ResizableBox</h1> -->
      <!--     <p>Left</p> -->
      <!--   </ResizableBox> -->
      <!--   <ResizableBox direction="right" class="w-full bg-green-700"> -->
      <!--     <h1>ResizableBox</h1> -->
      <!--     <p>Right</p> -->
      <!--   </ResizableBox> -->
      <!-- </div> -->
    </div>
  </div>
  <Settings/>
{/if}
