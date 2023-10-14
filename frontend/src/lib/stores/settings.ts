import {writable} from 'svelte/store'

export const showSettingsModal = writable(false)

export const PAGE_LINES_LOCAL_PATH = "settings:load:lines"
export const pageLinesLimit = writable(Number(localStorage.getItem(PAGE_LINES_LOCAL_PATH)) || 20)
