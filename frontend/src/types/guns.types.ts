import type { main } from "../../wailsjs/go/models";

export interface GetGuns {
  status: number
  data: Weapon[]
}

export interface Weapon {
  uuid: string
  displayName: string
  category: string
  defaultSkinUuid: string
  displayIcon: string
  killStreamIcon: string
  assetPath: string
  skins: Skin[]
}

export interface GridPosition {
  row: number
  column: number
}

export interface Skin {
  uuid: string
  displayName: string
  themeUuid: string
  contentTierUuid?: string
  displayIcon?: string
  wallpaper?: string
  assetPath: string
  chromas: Chroma[]
  levels: Level[]
}

export interface Chroma {
  uuid: string
  displayName: string
  displayIcon?: string
  fullRender: string
  swatch?: string
  streamedVideo?: string
  assetPath: string
}

export interface Level {
  uuid: string
  displayName: string
  levelItem?: string
  displayIcon?: string
  streamedVideo?: string
  assetPath: string
}


export type CustomGun = {
  id: string;
  skin: Skin;
  loadoutWeapon: main.Gun;
  weapon: Weapon
}
