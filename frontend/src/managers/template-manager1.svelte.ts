import type { Agent } from "src/types/agent.types";
import type { Weapon } from "src/types/guns.types";

export class ValorantManager {
    agents = $state<Agent[]>([])
    guns = $state<Weapon[]>([])

    async GetAgents() {
        const response = await fetch("https://valorant-api.com/v1/agents");
        const data = await response.json();
        this.agents = data.data;
    }

    async GetGuns() {
        const response = await fetch("https://valorant-api.com/v1/weapons");
        const data = await response.json();
        this.guns = data.data;
    }
}

let valorantManager: ValorantManager

export const useValorantManager = () => {
    if (!valorantManager) {
        valorantManager = new ValorantManager()
    }
    return valorantManager
}
