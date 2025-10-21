import type { main } from "wailsjs/go/models";

export type Template = {
    id: string;
    name: string;
    agent: string
    loadout: main.PlayerLoadoutResponse
};
