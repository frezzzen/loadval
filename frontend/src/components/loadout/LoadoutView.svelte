<script lang="ts">
    import { useValorantManager } from "../../managers/valorant-manager.svelte";
    import GunView from "./GunView.svelte";
    import { main } from "../../../wailsjs/go/models";
    import type {
        Chroma,
        CustomGun,
        Skin,
        Weapon,
    } from "../../types/guns.types";

    type Props = {
        loadout: main.PlayerLoadoutResponse;
    };

    let { loadout = $bindable() }: Props = $props();

    const valorantManager = useValorantManager();

    if (!valorantManager.guns.length) {
        valorantManager.GetGuns();
    }

    const guns = $derived.by(() => {
        return Object.values(valorantManager.guns).reduce(
            (acc: Record<string, CustomGun[]>, gun: Weapon) => {
                const category = gun.category;
                if (!acc[category]) {
                    acc[category] = [];
                }
                const loadoutGun = loadout.Guns.find(
                    (g: main.Gun) => g.ID === gun.uuid,
                );
                if (loadoutGun) {
                    const skin = gun.skins.find(
                        (s: Skin) => s.uuid === loadoutGun.SkinID,
                    );
                    if (skin) {
                        acc[category].push({
                            id: skin.uuid,
                            skin: skin,
                            loadoutWeapon: loadoutGun,
                            weapon: gun,
                        });
                    }
                }
                return acc;
            },
            {} as Record<string, CustomGun[]>,
        );
    });

    let selectedCustomGun = $state<CustomGun | undefined>(undefined);

    let currentView = $state<"loadout" | "gun-selection">("loadout");

    function handleSkinSelect(customGun: CustomGun) {
        selectedCustomGun = customGun;
        currentView = "gun-selection";
    }

    function handleBackToLoadout() {
        currentView = "loadout";
    }

    function handleSkinChange(
        skinUuid: string,
        chromaUuid?: string,
        levelUuid?: string,
    ) {
        if (!selectedCustomGun) return;

        const customGun = selectedCustomGun;

        const newSkin = customGun.weapon.skins.find(
            (s: Skin) => s.uuid === skinUuid,
        );
        if (newSkin) {
            const gunIndex = loadout.Guns.findIndex(
                (g: main.Gun) => g.ID === customGun.weapon.uuid,
            );
            if (gunIndex !== -1) {
                loadout.Guns[gunIndex].SkinID = skinUuid;

                if (chromaUuid) {
                    loadout.Guns[gunIndex].ChromaID = chromaUuid;
                } else if (newSkin.chromas && newSkin.chromas.length > 0) {
                    loadout.Guns[gunIndex].ChromaID = newSkin.chromas[0].uuid;
                }

                if (levelUuid) {
                    loadout.Guns[gunIndex].SkinLevelID = levelUuid;
                } else if (newSkin.levels && newSkin.levels.length > 0) {
                    loadout.Guns[gunIndex].SkinLevelID = newSkin.levels[0].uuid;
                }
            }

            selectedCustomGun = {
                ...customGun,
                skin: newSkin,
                id: newSkin.uuid,
                loadoutWeapon: main.Gun.createFrom({
                    ...customGun.loadoutWeapon,
                    SkinID: skinUuid,
                    ChromaID:
                        chromaUuid ||
                        newSkin.chromas?.[0]?.uuid ||
                        customGun.loadoutWeapon.ChromaID,
                    SkinLevelID:
                        levelUuid ||
                        newSkin.levels?.[0]?.uuid ||
                        customGun.loadoutWeapon.SkinLevelID,
                }),
            };
        }
    }

    const orderByCategory = {
        "EEquippableCategory::Heavy": {
            Ares: 0,
            Odin: 1,
        },
        "EEquippableCategory::Sniper": {
            Marshal: 0,
            Outlaw: 1,
            Operator: 2,
        },
        "EEquippableCategory::Rifle": {
            Bulldog: 0,
            Guardian: 1,
            Phantom: 2,
            Vandal: 3,
        },
        "EEquippableCategory::SMG": {
            Stinger: 0,
            Spectre: 1,
        },
        "EEquippableCategory::Shotgun": {
            Bucky: 0,
            Judge: 1,
        },
        "EEquippableCategory::Sidearm": {
            Classic: 0,
            Shorty: 1,
            Frenzy: 2,
            Ghost: 3,
            Sheriff: 4,
        },
        "EEquippableCategory::Melee": {
            Melee: 0,
        },
    };

    function getSortedGuns(category: string) {
        return guns[category].sort((a, b) => {
            console.log(a.weapon.displayName, b.weapon.displayName);
            const aIndex = orderByCategory[category][a.weapon.displayName];
            const bIndex = orderByCategory[category][b.weapon.displayName];
            return aIndex - bIndex;
        });
    }
</script>

{#snippet categoryColumn(category: string)}
    <div class="container">
        <div class="row">
            {#each getSortedGuns(category) as customGun (customGun.id)}
                {@const chroma = customGun.skin.chromas.find(
                    (c: Chroma) => c.uuid === customGun.loadoutWeapon.ChromaID,
                )}
                {@const displayIcon =
                    chroma?.displayIcon ||
                    customGun.skin.displayIcon ||
                    customGun.weapon.displayIcon}
                <div
                    class="row-item"
                    onclick={() => handleSkinSelect(customGun)}
                    onkeydown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                            handleSkinSelect(customGun);
                        }
                    }}
                    role="button"
                    tabindex="0"
                >
                    <img src={displayIcon} alt={customGun.skin.displayName} />
                </div>
            {/each}
        </div>
    </div>
{/snippet}

{#if currentView === "loadout"}
    <div class="loadout-view">
        {#if Object.keys(guns).length > 0}
            <div class="column">
                {@render categoryColumn("EEquippableCategory::Sidearm")}
            </div>
            <div class="column">
                {@render categoryColumn("EEquippableCategory::SMG")}
                {@render categoryColumn("EEquippableCategory::Shotgun")}
            </div>
            <div class="column">
                {@render categoryColumn("EEquippableCategory::Rifle")}
                {@render categoryColumn("EEquippableCategory::Melee")}
            </div>
            <div class="column">
                {@render categoryColumn("EEquippableCategory::Sniper")}
                {@render categoryColumn("EEquippableCategory::Heavy")}
            </div>
        {:else}
            <div class="loading">Loading...</div>
        {/if}
    </div>
{:else if selectedCustomGun}
    <GunView
        customGun={selectedCustomGun}
        onBack={handleBackToLoadout}
        onSelect={handleSkinChange}
    />
{/if}

<style lang="scss">
    .loadout-view {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 1fr;
        gap: 4rem;
    }

    .column {
        display: flex;
        flex-direction: column;
        height: 100%;
        gap: 1rem;
        width: 100%;

        .column-title {
            font-size: 1.4rem;
            font-weight: 600;
            color: #cccaca;
        }

        .row {
            display: flex;
            flex-direction: column;
            gap: 1rem;
        }

        .row-item {
            height: 14rem;
            border-radius: 1rem;
            overflow: hidden;
            border: 2px solid rgba(60, 60, 60, 0.8);
            cursor: pointer;
            transition: all 0.2s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.03),
                rgba(122, 40, 203, 0.01)
            );
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
            position: relative;

            &::before {
                content: "";
                position: absolute;
                inset: 0;
                border-radius: 1rem;
                padding: 2px;
                background: linear-gradient(
                    135deg,
                    transparent,
                    rgba(173, 64, 255, 0)
                );
                -webkit-mask:
                    linear-gradient(#fff 0 0) content-box,
                    linear-gradient(#fff 0 0);
                -webkit-mask-composite: xor;
                mask:
                    linear-gradient(#fff 0 0) content-box,
                    linear-gradient(#fff 0 0);
                mask-composite: exclude;
                opacity: 0;
                transition: opacity 0.2s ease;
            }

            &:hover {
                border-color: rgba(173, 64, 255, 0.6);
                background: linear-gradient(
                    135deg,
                    rgba(173, 64, 255, 0.08),
                    rgba(122, 40, 203, 0.04)
                );
                transform: translateY(-4px);
                box-shadow: 0 8px 24px rgba(173, 64, 255, 0.25);

                &::before {
                    opacity: 1;
                }

                img {
                    transform: scale(1.05);
                    filter: drop-shadow(0 0 20px rgba(173, 64, 255, 0.4));
                }
            }

            &:active {
                transform: translateY(-2px);
                box-shadow: 0 4px 16px rgba(173, 64, 255, 0.2);
            }

            &:focus-visible {
                outline: none;
                border-color: rgba(173, 64, 255, 0.8);
                box-shadow: 0 0 0 3px rgba(173, 64, 255, 0.2);
            }

            img {
                width: 70%;
                height: 70%;
                object-fit: contain;
                transition: all 0.3s ease;
                z-index: 1;
            }
        }
    }
</style>
