import { writable } from "svelte/store";
import type {src} from "../../../wailsjs/go/models"

export const imagesArrayData = writable<src.Image[]>([])
