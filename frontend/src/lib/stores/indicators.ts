import { writable } from "svelte/store";

export const loadingFileIndicator = writable(false)

export const cacheLoading = writable(false)

export const saveImagesIndicator = writable(false)