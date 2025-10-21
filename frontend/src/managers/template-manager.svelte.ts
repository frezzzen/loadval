import { SetPlayerLoadout } from "../../wailsjs/go/main/ValorantAPI"
import type { Template } from "../types/template.type"
import type { main } from "../../wailsjs/go/models"

export class TemplateManager {
    templates = $state<Template[]>([])
    agentLoadouts = $state<Map<string, main.PlayerLoadoutResponse>>(new Map())
    isAgentLoadoutsEnabled = $state<boolean>(localStorage.getItem("isAgentLoadoutsEnabled") === "true")
    activeTemplate = $state<string | null>(localStorage.getItem("activeTemplate") || null)

    constructor() {
        const templates = localStorage.getItem("templates")
        if (templates) {
            this.templates = JSON.parse(templates)
        }

        const agentLoadouts = localStorage.getItem("agentLoadouts")
        if (agentLoadouts) {
            const parsed = JSON.parse(agentLoadouts)
            this.agentLoadouts = new Map(Object.entries(parsed))
        }
    }

    setIsAgentLoadoutsEnabled(enabled: boolean) {
        this.isAgentLoadoutsEnabled = enabled
        localStorage.setItem("isAgentLoadoutsEnabled", enabled.toString())
    }

    setActiveTemplate(template: string | null) {
        this.activeTemplate = template
        localStorage.setItem("activeTemplate", template)
    }

    getActiveTemplate() {
        return this.activeTemplate
    }

    addTemplate(template: Template) {
        this.templates.push(template)
        localStorage.setItem("templates", JSON.stringify(this.templates))
        return template
    }

    removeTemplate(id: string) {
        this.templates = this.templates.filter((template) => template.id !== id)
        localStorage.setItem("templates", JSON.stringify(this.templates))
    }

    getTemplate(id: string) {
        return this.templates.find((template) => template.id === id)
    }

    updateTemplate(id: string, template: Template) {
        this.templates = this.templates.map((t) => t.id === id ? template : t)
        localStorage.setItem("templates", JSON.stringify(this.templates))
    }

    getTemplates() {
        return this.templates
    }

    setAgentLoadout(agentUuid: string, loadout: main.PlayerLoadoutResponse) {
        this.agentLoadouts.set(agentUuid, loadout)
        const obj = Object.fromEntries(this.agentLoadouts)
        localStorage.setItem("agentLoadouts", JSON.stringify(obj))
    }

    getAgentLoadout(agentUuid: string) {
        return this.agentLoadouts.get(agentUuid)
    }

    hasAgentLoadout(agentUuid: string) {
        return this.agentLoadouts.has(agentUuid)
    }

    removeAgentLoadout(agentUuid: string) {
        this.agentLoadouts.delete(agentUuid)
        const obj = Object.fromEntries(this.agentLoadouts)
        localStorage.setItem("agentLoadouts", JSON.stringify(obj))
    }


    async saveLoadout(loadout: main.PlayerLoadoutResponse) {
        try {
            const result = await SetPlayerLoadout(loadout);
            loadout = result;
        } catch (error) {
            console.error(error);
        }
    }
}



let templateManager: TemplateManager

export const useTemplateManager = () => {
    if (!templateManager) {
        templateManager = new TemplateManager()
    }
    return templateManager
}
