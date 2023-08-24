<script lang="ts">
  import {writable} from 'svelte/store'
    import { fade, scale, slide } from 'svelte/transition';
    export let show = writable(false);
    export let title = 'Modal Title'
</script>
{#if $show}
    <div class="modal w-screen opacity-100 visible" transition:scale={{start: 3}}>
        <span class="modal-overlay"></span>
        <div class="modal-content w-full lg:max-w-2xl flex flex-col gap-5">
            <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" on:click={() => show.set(false)}>✕</button>
            <h2 class="text-xl">{title}</h2>
                <slot></slot>
        </div>
    </div>
{/if}
