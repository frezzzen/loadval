import { SetPlayerLoadout } from "../../wailsjs/go/main/ValorantAPI"
import {
    LoadTemplates,
    SaveTemplates,
    LoadAgentLoadouts,
    SaveAgentLoadouts,
    LoadSettings,
    SaveSettings
} from "../../wailsjs/go/main/StorageAPI"
import type { Template } from "../types/template.type"
import type { main } from "../../wailsjs/go/models"

export class TemplateManager {
    templates = $state<Template[]>([])
    agentLoadouts = $state<Map<string, main.PlayerLoadoutResponse>>(new Map())
    isAgentLoadoutsEnabled = $state<boolean>(false)
    activeTemplate = $state<string | null>(null)

    constructor() {
        this.loadData()
    }

    async loadData() {
        try {
            const templatesData = await LoadTemplates()
            if (templatesData && templatesData !== "{}") {
                this.templates = JSON.parse(templatesData)
            }

            const agentLoadoutsData = await LoadAgentLoadouts()
            if (agentLoadoutsData && agentLoadoutsData !== "{}") {
                const parsed = JSON.parse(agentLoadoutsData)
                this.agentLoadouts = new Map(Object.entries(parsed))
            }

            const settingsData = await LoadSettings()
            if (settingsData && settingsData !== "{}") {
                const settings = JSON.parse(settingsData)
                this.isAgentLoadoutsEnabled = settings.isAgentLoadoutsEnabled || false
                this.activeTemplate = settings.activeTemplate || null
            }
        } catch (error) {
            console.error("Failed to load data from storage:", error)
            this.migrateFromLocalStorage()
        }
    }

    migrateFromLocalStorage() {
        const templates = localStorage.getItem("templates")
        if (templates) {
            this.templates = JSON.parse(templates)
            this.saveTemplates()
        }

        const agentLoadouts = localStorage.getItem("agentLoadouts")
        if (agentLoadouts) {
            const parsed = JSON.parse(agentLoadouts)
            this.agentLoadouts = new Map(Object.entries(parsed))
            this.saveAgentLoadouts()
        }

        const isEnabled = localStorage.getItem("isAgentLoadoutsEnabled")
        if (isEnabled) {
            this.isAgentLoadoutsEnabled = isEnabled === "true"
        }

        const activeTemplate = localStorage.getItem("activeTemplate")
        if (activeTemplate) {
            this.activeTemplate = activeTemplate
        }

        this.saveSettings()

        localStorage.removeItem("templates")
        localStorage.removeItem("agentLoadouts")
        localStorage.removeItem("isAgentLoadoutsEnabled")
        localStorage.removeItem("activeTemplate")
    }

    async saveTemplates() {
        try {
            await SaveTemplates(JSON.stringify(this.templates))
        } catch (error) {
            console.error("Failed to save templates:", error)
        }
    }

    async saveAgentLoadouts() {
        try {
            const obj = Object.fromEntries(this.agentLoadouts)
            await SaveAgentLoadouts(JSON.stringify(obj))
        } catch (error) {
            console.error("Failed to save agent loadouts:", error)
        }
    }

    async saveSettings() {
        try {
            const settings = {
                isAgentLoadoutsEnabled: this.isAgentLoadoutsEnabled,
                activeTemplate: this.activeTemplate
            }
            await SaveSettings(JSON.stringify(settings))
        } catch (error) {
            console.error("Failed to save settings:", error)
        }
    }

    setIsAgentLoadoutsEnabled(enabled: boolean) {
        this.isAgentLoadoutsEnabled = enabled
        this.saveSettings()
    }

    setActiveTemplate(template: string | null) {
        this.activeTemplate = template
        this.saveSettings()
    }

    getActiveTemplate() {
        return this.activeTemplate
    }

    addTemplate(template: Template) {
        this.templates.push(template)
        this.saveTemplates()
        return template
    }

    removeTemplate(id: string) {
        this.templates = this.templates.filter((template) => template.id !== id)
        this.saveTemplates()
    }

    getTemplate(id: string) {
        return this.templates.find((template) => template.id === id)
    }

    updateTemplate(id: string, template: Template) {
        this.templates = this.templates.map((t) => t.id === id ? template : t)
        this.saveTemplates()
    }

    getTemplates() {
        return this.templates
    }

    setAgentLoadout(agentUuid: string, loadout: main.PlayerLoadoutResponse) {
        const clonedLoadout = JSON.parse(JSON.stringify(loadout))
        this.agentLoadouts.set(agentUuid, clonedLoadout)
        this.saveAgentLoadouts()
    }

    getAgentLoadout(agentUuid: string) {
        const loadout = this.agentLoadouts.get(agentUuid)
        return loadout ? JSON.parse(JSON.stringify(loadout)) : undefined
    }

    hasAgentLoadout(agentUuid: string) {
        return this.agentLoadouts.has(agentUuid)
    }

    removeAgentLoadout(agentUuid: string) {
        this.agentLoadouts.delete(agentUuid)
        this.saveAgentLoadouts()
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
