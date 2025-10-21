import type { main } from "../../wailsjs/go/models";

export class LoadoutManager {
    loadout = $state<main.PlayerLoadoutResponse | null>(null);
    ownedItems = $state<main.OwnedItemsResponseEntitlement[] | null>(null);





    getLoadout() {
        return this.loadout;
    }


    getOwnedItems() {
        return this.ownedItems;
    }


    setLoadout(loadout: main.PlayerLoadoutResponse) {
        this.loadout = loadout;
    }

    setOwnedItems(ownedItems: main.OwnedItemsResponseEntitlement[]) {
        this.ownedItems = ownedItems;
    }
}

let loadoutManager: LoadoutManager

export const useLoadoutManager = () => {
    if (!loadoutManager) {
        loadoutManager = new LoadoutManager()
    }
    return loadoutManager
}
