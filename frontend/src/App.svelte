<script lang="ts">
    import Header from './lib/components/Header.svelte'
    import Settings from "./lib/settings/Settings.svelte";
    import {cacheLoading, loadingFileIndicator, saveImagesIndicator} from './lib/stores/indicators'
    import {onMount} from 'svelte';
    import type {src} from '../wailsjs/go/models'
    import {EventsOn} from '../wailsjs/runtime'
    import autoAnimate from '@formkit/auto-animate'
    import {fade} from 'svelte/transition'
    import Images from './lib/home/Images.svelte';
    import Categories from './lib/home/Categories.svelte';
    import Tags from './lib/home/Tags.svelte';

    let activeTab = 0
    const activeTabClasses = ' tab-underline tab-active '


    let cacheStatusText = ""
    let saveImagesStatusText = ""
    onMount(() => {
        EventsOn("cache-start", () => {
            cacheLoading.set(true)
        })

        EventsOn("cache-end", () => {
            cacheLoading.set(false)
        })
        EventsOn("cache-event", (d: { image: src.Image, totalImages: number, cachedImages: number }) => {
            cacheStatusText = `Caching Image ${d.image.id}. (${d.cachedImages}/${d.totalImages})`
        })

        EventsOn("save-images-start", () => {
            saveImagesIndicator.set(true)
        })

        EventsOn("save-images-end", () => {
            saveImagesIndicator.set(false)
        })
        EventsOn("image-saved", (d: {image: any, totalImages: number, savedImages: number}) => {
            saveImagesStatusText = `Saved ${d.image.id} (${d.savedImages}/${d.totalImages})`
            console.log(d.image)
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
{:else if $saveImagesIndicator}
    <div class="w-full h-screen" use:autoAnimate>
        <div class="flex flex-col p-10 text-center w-full h-full items-center justify-center">
            <h1 class="text-2xl">
                Saving Images...
            </h1>
            <p class="my-5">
                {saveImagesStatusText}
            </p>
        </div>
    </div>
{:else}
    <Header/>
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
                <div class="tabs p-3 w-full justify-center bg-slate-3">
                    <button on:click={() => activeTab = 0}
                            class="{activeTab === 0 || activeTab > 2 ? activeTabClasses : ''} tab px-6 transition duration-200">
                        Images
                    </button>
                    <button on:click={() => activeTab = 1}
                            class="{activeTab === 1 ? activeTabClasses : '' } transition duration-200 tab px-6">
                        Categories
                    </button>
                    <button on:click={() => activeTab = 2}
                            class="{activeTab === 2 ? activeTabClasses : '' } transition duration-200 tab px-6">Tags
                    </button>
                </div>
                <div class="tabs-content">
                    {#if activeTab === 1}
                        <div in:fade>
                            <Categories/>
                        </div>
                    {:else if activeTab === 2}
                        <div in:fade>
                            <Tags/>
                        </div>
                    {:else}
                        <div class="flex flex-col" style="gap: 10px" in:fade>
                            <Images/>
                        </div>
                    {/if}
                </div>
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
