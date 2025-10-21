<script lang="ts">
    import { onMount } from "svelte";

    import {
        GetMainData,
        GetPreGameMatch,
        GetPreGamePlayer,
        SetPlayerLoadout,
    } from "../../wailsjs/go/main/ValorantAPI";
    import type { main } from "../../wailsjs/go/models";
    import TemplateMenu from "../components/templates/TemplateMenu.svelte";
    import { useTemplateManager } from "../managers/template-manager.svelte";
    import { useLoadoutManager } from "../managers/loadout-manager.svelte";

    type Props = {
        loadout: main.PlayerLoadoutResponse;
        ownedItems: main.OwnedItemsResponseEntitlement[];
    };

    let { loadout, ownedItems }: Props = $props();

    const templateManager = useTemplateManager();
    const loadoutManager = useLoadoutManager();

    onMount(async () => {
        loadoutManager.setLoadout(loadout);
        loadoutManager.setOwnedItems(ownedItems);
    });

    async function saveLoadout(loadout: main.PlayerLoadoutResponse) {
        console.log(loadout);
        try {
            const result = await SetPlayerLoadout(loadout);
            loadout = result;
            console.log("Loadout saved");
        } catch (error) {
            console.error(error);
        }
    }

    async function run() {
        let interval = setInterval(async () => {
            const preGamePlayer = await GetPreGamePlayer();
            if (preGamePlayer) {
                const me = preGamePlayer.Subject;
                let secondInterval = setInterval(async () => {
                    clearInterval(interval);
                    const preGameMatch = await GetPreGameMatch();
                    if (preGameMatch) {
                        console.log(preGameMatch);
                        const myTeam = preGameMatch.Teams.find((team) =>
                            team.Players.some(
                                (player) => player.Subject === me,
                            ),
                        );
                        const mePlayer = myTeam?.Players.find(
                            (player) => player.Subject === me,
                        );
                        if (mePlayer) {
                            console.log(mePlayer.CharacterSelectionState);
                            if (mePlayer.CharacterSelectionState === "locked") {
                                console.log(
                                    "locked",
                                    templateManager.agentLoadouts,
                                );
                                const agentLoadout =
                                    templateManager.getAgentLoadout(
                                        mePlayer.CharacterID,
                                    );
                                if (agentLoadout) {
                                    console.log(agentLoadout);
                                    await saveLoadout(agentLoadout);
                                    clearInterval(secondInterval);
                                    setTimeout(() => {
                                        run();
                                    }, 5000);
                                }
                            }
                        }
                    }
                }, 1000);
            }
        }, 5000);
    }

    run();
</script>

<TemplateMenu {loadout} />

<style lang="scss">
    .title {
        font-size: 1.6rem;
        color: white;
        font-weight: 600;
        text-align: center;
        margin-bottom: 1rem;
        text-align: center;
    }

    drag {
        position: fixed;
        top: 0;
        right: 0;
        z-index: 1000;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 1rem;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.1),
            rgba(255, 255, 255, 0.05)
        );
        backdrop-filter: blur(10px);
        border-bottom-left-radius: 1rem;
        border: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
        cursor: move;
        transition: all 0.3s ease;

        &:hover {
            background: linear-gradient(
                135deg,
                rgba(255, 255, 255, 0.15),
                rgba(255, 255, 255, 0.08)
            );
            box-shadow: 0 6px 25px rgba(0, 0, 0, 0.4);
        }

        button {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 3.2rem;
            height: 3.2rem;
            border: none;
            border-radius: 0.8rem;
            background: rgba(255, 255, 255, 0.1);
            color: #cccaca;
            cursor: pointer;
            transition: all 0.2s ease;
            backdrop-filter: blur(5px);
            border: 1px solid rgba(255, 255, 255, 0.1);

            &:hover {
                background: rgba(255, 255, 255, 0.2);
                transform: translateY(-1px);
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
            }

            &:active {
                transform: translateY(0);
                box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
            }

            &:nth-child(1) {
                &:hover {
                    background: rgba(34, 197, 94, 0.2);
                    color: #22c55e;
                }
            }

            &:nth-child(2) {
                &:hover {
                    background: rgba(251, 191, 36, 0.2);
                    color: #fbbf24;
                }
            }

            &:nth-child(3) {
                &:hover {
                    background: rgba(239, 68, 68, 0.2);
                    color: #ef4444;
                }
            }

            .material-icons {
                font-size: 1.8rem;
                transition: transform 0.2s ease;
            }

            &:hover .material-icons {
                transform: scale(1.1);
            }
        }
    }

    @media (max-width: 768px) {
        drag {
            padding: 0.8rem;
            gap: 0.4rem;

            button {
                width: 2.8rem;
                height: 2.8rem;

                .material-icons {
                    font-size: 1.6rem;
                }
            }
        }
    }

    @media (max-width: 480px) {
        drag {
            padding: 0.6rem;
            gap: 0.3rem;

            button {
                width: 2.4rem;
                height: 2.4rem;

                .material-icons {
                    font-size: 1.4rem;
                }
            }
        }
    }
</style>
