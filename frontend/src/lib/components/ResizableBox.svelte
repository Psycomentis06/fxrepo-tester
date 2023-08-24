<script lang="ts">
    import { onMount } from "svelte";

  type DirectionAttrs = {
    cursor: string
    bar: string
    container: string
  }
  type ActiveDirection = {
    right: DirectionAttrs
    left: DirectionAttrs
    top: DirectionAttrs
    bottom: DirectionAttrs
  }
  export let direction: 'left' | 'right' | 'top' | 'bottom' = 'top';

  const directionClasses: ActiveDirection = {
    left: {
      cursor: 'cursor-ew-resize',
      bar: 'w-8',
      container: 'flex-row',
    },
    right: {
      cursor: 'cursor-ew-resize',
      bar: 'w-8',
      container: 'flex-row-reverse'
    },
    top: {
      cursor: 'cursor-ns-resize',
      bar: 'w-full h-8',
      container: 'flex-col'
    },
    bottom: {
      cursor: 'cursor-ns-resize',
      bar: 'w-full h-8',
      container: 'flex-col-reverse'
    },
  }
  let activeDirection = directionClasses[direction] ;

  var contentElement: HTMLDivElement
  let widthInitSize: number
  var widthInitPos:  number
  let heightInitSize: number
  var heightInitPos: number
  var isDragging = false

    const mouseDown = (e: MouseEvent) => {
      isDragging = true
      widthInitSize = parseInt(window.getComputedStyle(contentElement).width)
      widthInitPos = e.clientX
      heightInitSize = parseInt(window.getComputedStyle(contentElement).height)
      heightInitPos = e.clientY
    }

    const mouseUp = (e: MouseEvent) => {
      isDragging = false
    }

    const mouseMove = (e: MouseEvent) => {
      if (!isDragging) return
      if(direction === 'right' || direction === 'left') {
        contentElement.style.width = widthInitSize + (e.clientX - widthInitPos) + 'px'
      } else if(direction === 'top' || direction === 'bottom') {
        contentElement.style.height = heightInitSize + (e.clientY - heightInitPos) + 'px'
      }
    }
</script>

<div class="resizable-box {$$props.class}" {...$$restProps} bind:this={contentElement}>
  <div class="flex {activeDirection.container}">
    <div class="{($$props['bar-class'] || '')  + ' ' + activeDirection.cursor + ' ' + activeDirection.bar} z-[1000] box-bar bg-gray-800 flex items-center justify-center cursor-ew-resize" style="gap: 2px;" on:mousedown={mouseDown} on:mouseup={mouseUp} on:mousemove={mouseMove} on:mouseleave={mouseUp}>
      <div class="pointer-events-auto w-2 h-2 rounded-full bg-gray-600"></div>
      <div class="pointer-events-auto w-2 h-2 rounded-full bg-gray-600"></div>
      <div class="pointer-events-auto w-2 h-2 rounded-full bg-gray-600"></div>
    </div>
    <div class="w-full">
      <slot />
    </div>
  </div>
</div>
