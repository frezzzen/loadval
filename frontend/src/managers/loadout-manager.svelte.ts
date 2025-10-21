import type { main } from "../../wailsjs/go/models";

export class LoadoutManager {
    loadout = $state<main.PlayerLoadoutResponse | null>(null);
    ownedSkins = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
    ownedSkinVariants = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
    ownedAgents = $state<main.OwnedItemsResponseEntitlement[] | null>(null);
    ownedCards = $state<main.OwnedItemsResponseEntitlement[] | null>(null);





    getLoadout() {
        return this.loadout;
    }


    getOwnedSkins() {
        return this.ownedSkins;
    }
    getOwnedSkinVariants() {
        return this.ownedSkinVariants;
    }
    getOwnedAgents() {
        return this.ownedAgents;
    }
    getOwnedCards() {
        return this.ownedCards;
    }


    setLoadout(loadout: main.PlayerLoadoutResponse) {
        this.loadout = loadout;
    }

    setOwnedSkins(ownedSkins: main.OwnedItemsResponseEntitlement[]) {
        this.ownedSkins = ownedSkins;
    }
    setOwnedSkinVariants(ownedSkinVariants: main.OwnedItemsResponseEntitlement[]) {
        this.ownedSkinVariants = ownedSkinVariants;
    }
    setOwnedAgents(ownedAgents: main.OwnedItemsResponseEntitlement[]) {
        this.ownedAgents = ownedAgents;
    }
    setOwnedCards(ownedCards: main.OwnedItemsResponseEntitlement[]) {
        this.ownedCards = ownedCards;
    }
}

let loadoutManager: LoadoutManager

export const useLoadoutManager = () => {
    if (!loadoutManager) {
        loadoutManager = new LoadoutManager()
    }
    return loadoutManager
}
