<script lang="ts">
    import { useLoadoutManager } from "../../managers/loadout-manager.svelte";
    import type { CustomGun, Skin } from "../../types/guns.types";
    import CoreSwiper from "../core/core-swiper.svelte";

    type Props = {
        onSelect?: (skin: string, chroma?: string, level?: string) => void;
        onBack?: () => void;
        customGun: CustomGun;
    };

    let { customGun, onSelect, onBack }: Props = $props();

    const loadoutManager = useLoadoutManager();

    let selectedSkin = $state<string>(customGun.skin.uuid);
    let selectedChroma = $state<string>(customGun.loadoutWeapon.ChromaID);
    let selectedLevel = $state<string | undefined>(
        customGun.loadoutWeapon.SkinLevelID,
    );
    let hoveredSkinImage = $state<string | null>(null);

    const currentSkin = $derived.by(() => {
        return customGun.weapon.skins.find((s) => s.uuid === selectedSkin);
    });

    const selectedSkinImage = $derived.by(() => {
        if (hoveredSkinImage) return hoveredSkinImage;

        const skin = currentSkin;
        if (!skin) return customGun.skin.displayIcon || "";

        // Try to find the selected chroma
        const chroma = skin.chromas.find((c) => c.uuid === selectedChroma);
        return (
            chroma?.displayIcon ||
            skin.displayIcon ||
            customGun.weapon.displayIcon
        );
    });

    const ownedSkins = $derived.by((): Skin[] => {
        return customGun.weapon.skins.filter((skin) => {
            const isInSkins = loadoutManager.ownedItems?.some(
                (ownedItem) => ownedItem.ItemID === skin.uuid,
            );
            if (isInSkins) {
                return true;
            }
            const levels = skin.levels.filter((level) => {
                const isInLevels = loadoutManager.ownedItems?.some(
                    (ownedItem) => ownedItem.ItemID === level.uuid,
                );
                return isInLevels;
            });
            return levels.length > 0;
        });
    });

    $inspect(ownedSkins);

    function handleBackClick() {
        if (onBack) {
            onBack();
        }
    }

    function handleSkinSelect(skinUuid: string) {
        selectedSkin = skinUuid;
        const skin = customGun.weapon.skins.find((s) => s.uuid === skinUuid);
        if (skin) {
            // Auto-select first chroma and level
            if (skin.chromas && skin.chromas.length > 0) {
                selectedChroma = skin.chromas[0].uuid;
            }
            if (skin.levels && skin.levels.length > 0) {
                selectedLevel = skin.levels[skin.levels.length - 1].uuid;
            }
        }
        if (onSelect) {
            onSelect(skinUuid, selectedChroma, selectedLevel);
        }
    }

    function handleChromaSelect(chromaUuid: string) {
        selectedChroma = chromaUuid;
        if (onSelect) {
            onSelect(selectedSkin, chromaUuid, selectedLevel);
        }
    }

    function handleLevelSelect(levelUuid: string) {
        selectedLevel = levelUuid;
        if (onSelect) {
            onSelect(selectedSkin, selectedChroma, levelUuid);
        }
    }

    $effect(() => {
        selectedSkin = customGun.skin.uuid;
        selectedChroma = customGun.loadoutWeapon.ChromaID;
        selectedLevel = customGun.loadoutWeapon.SkinLevelID;
    });

    $inspect(ownedSkins);
</script>

<div class="gun-view">
    <div class="header">
        <button class="back-button" onclick={handleBackClick}>
            <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    d="M19 12H5M5 12L12 19M5 12L12 5"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
            Back to Loadout
        </button>
        <h2 class="title">Select Skin</h2>
    </div>

    <div class="hero">
        <img src={selectedSkinImage} alt="Selected Gun" />
        <div class="skin-info">
            <h3>{currentSkin?.displayName || customGun.skin.displayName}</h3>
        </div>
        {#if currentSkin && currentSkin.chromas && currentSkin.chromas.length > 1}
            <div class="customization-section chromas">
                <div class="chromas">
                    {#each currentSkin.chromas as chroma (chroma.uuid)}
                        <div
                            class="chroma-item"
                            class:selected={chroma.uuid === selectedChroma}
                            role="button"
                            tabindex="0"
                            onclick={() => handleChromaSelect(chroma.uuid)}
                            onkeydown={(e) => {
                                if (e.key === "Enter" || e.key === " ") {
                                    handleChromaSelect(chroma.uuid);
                                }
                            }}
                            onmouseenter={() => {
                                hoveredSkinImage =
                                    chroma.displayIcon || chroma.fullRender;
                            }}
                            onmouseleave={() => {
                                hoveredSkinImage = null;
                            }}
                        >
                            {#if chroma.swatch}
                                <div
                                    class="chroma-swatch"
                                    style="background-image: url({chroma.swatch})"
                                ></div>
                            {:else if chroma.displayIcon}
                                <img
                                    src={chroma.displayIcon}
                                    alt={chroma.displayName}
                                />
                            {/if}
                            {#if chroma.uuid === selectedChroma}
                                <div class="selected-badge">
                                    <svg
                                        width="12"
                                        height="12"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path
                                            d="M20 6L9 17L4 12"
                                            stroke="currentColor"
                                            stroke-width="3"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                    </svg>
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
        {#if currentSkin && currentSkin.levels && currentSkin.levels.length > 1}
            <div class="customization-section levels">
                <div class="levels">
                    {#each currentSkin.levels as level (level.uuid)}
                        <div
                            class="level-item"
                            class:selected={level.uuid === selectedLevel}
                            role="button"
                            tabindex="0"
                            onclick={() => handleLevelSelect(level.uuid)}
                            onkeydown={(e) => {
                                if (e.key === "Enter" || e.key === " ") {
                                    handleLevelSelect(level.uuid);
                                }
                            }}
                        >
                            <div class="level-icon">
                                <span class="level-number"
                                    >{level.displayName.match(/\d+/)?.[0] ||
                                        "1"}</span
                                >
                            </div>
                            <div class="level-name">{level.displayName}</div>
                            {#if level.uuid === selectedLevel}
                                <div class="selected-badge">
                                    <svg
                                        width="12"
                                        height="12"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path
                                            d="M20 6L9 17L4 12"
                                            stroke="currentColor"
                                            stroke-width="3"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                        />
                                    </svg>
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>

    <div class="customization-section">
        <h3 class="section-title">Skins</h3>
        <div class="other-guns">
            <CoreSwiper options={{ slidesPerView: 8, spaceBetween: 10 }}>
                {#each ownedSkins as skin (skin.uuid)}
                    <div
                        class="other-gun swiper-slide"
                        class:selected={skin.uuid === selectedSkin}
                        role="button"
                        tabindex="0"
                        onclick={() => handleSkinSelect(skin.uuid)}
                        onkeydown={(e) => {
                            if (e.key === "Enter" || e.key === " ") {
                                handleSkinSelect(skin.uuid);
                            }
                        }}
                    >
                        <img src={skin.displayIcon} alt={skin.displayName} />
                        <div class="gun-name">{skin.displayName}</div>
                        {#if skin.uuid === selectedSkin}
                            <div class="selected-badge">
                                <svg
                                    width="16"
                                    height="16"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        d="M20 6L9 17L4 12"
                                        stroke="currentColor"
                                        stroke-width="3"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                    />
                                </svg>
                            </div>
                        {/if}
                    </div>
                {/each}
            </CoreSwiper>
        </div>
    </div>
</div>

<style lang="scss">
    .gun-view {
        display: flex;
        flex-direction: column;
        gap: 3rem;
        height: 100%;
        width: 100%;

        .header {
            display: flex;
            align-items: center;
            gap: 2rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid #3c3c3c;
            width: 100%;

            .back-button {
                display: flex;
                align-items: center;
                gap: 0.5rem;
                padding: 1rem 2rem;
                background: rgba(255, 255, 255, 0.1);
                border: 1px solid rgba(255, 255, 255, 0.2);
                border-radius: 10px;
                color: #ffffff;
                font-size: 1.05rem;
                font-weight: 600;
                cursor: pointer;
                transition: all 0.2s ease;
                backdrop-filter: blur(10px);
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

                &:hover {
                    background: rgba(255, 255, 255, 0.15);
                    border-color: rgba(255, 255, 255, 0.3);
                    transform: translateY(-2px);
                    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
                }

                &:active {
                    transform: translateY(0);
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                }

                svg {
                    width: 20px;
                    height: 20px;
                }
            }

            .title {
                font-size: 2rem;
                font-weight: 600;
                color: #ffffff;
                margin: 0;
            }
        }

        .hero {
            width: 100%;
            height: 40rem;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            background: linear-gradient(
                135deg,
                rgba(173, 64, 255, 0.1) 0%,
                rgba(60, 60, 60, 0.1) 100%
            );
            border-radius: 1rem;
            border: 2px solid #3c3c3c;
            padding: 2rem;
            position: relative;
            overflow: visible;

            &::before {
                content: "";
                position: absolute;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background: radial-gradient(
                    circle at 50% 50%,
                    rgba(173, 64, 255, 0.2) 0%,
                    transparent 70%
                );
                pointer-events: none;
            }

            img {
                width: 80%;
                height: 80%;
                object-fit: contain;
                position: relative;
                z-index: 1;
                filter: drop-shadow(0 10px 30px rgba(0, 0, 0, 0.5));
            }

            .skin-info {
                position: absolute;
                top: 1.5rem;
                left: 50%;
                transform: translateX(-50%);
                z-index: 2;
                background: rgba(0, 0, 0, 0.7);
                padding: 0.75rem 1.5rem;
                border-radius: 0.5rem;
                backdrop-filter: blur(10px);

                h3 {
                    margin: 0;
                    font-size: 1.25rem;
                    font-weight: 600;
                    color: #ffffff;
                    text-align: center;
                }
            }

            .customization-section {
                position: absolute;
                z-index: 2;
            }

            .levels {
                bottom: 1.5rem;
                right: 1.5rem;
            }

            .chromas {
                bottom: 1.5rem;
                left: 1.5rem;
            }
        }

        .customization-section {
            display: flex;
            flex-direction: column;
            gap: 1rem;

            &:not(.hero .customization-section) {
                .section-title {
                    font-size: 1.5rem;
                    font-weight: 600;
                    color: #ffffff;
                    margin: 0;
                    padding-bottom: 0.5rem;
                    border-bottom: 1px solid #3c3c3c;
                }
            }
        }

        .other-guns {
            padding: 0.5rem;
        }

        .other-gun {
            position: relative;
            height: 14rem;
            border-radius: 0.75rem;
            overflow: hidden;
            border: 2px solid #3c3c3c;
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            background: rgba(30, 30, 30, 0.5);
            padding: 1rem;

            &:hover {
                border-color: #ad40ff;
                transform: translateY(-4px);
                box-shadow: 0 8px 20px rgba(173, 64, 255, 0.3);

                img {
                    transform: scale(1.15);
                }
            }

            &.selected {
                border-color: #ad40ff;
                background: rgba(173, 64, 255, 0.15);
                box-shadow: 0 0 20px rgba(173, 64, 255, 0.4);
            }

            &:focus {
                outline: 2px solid #ad40ff;
                outline-offset: 2px;
            }

            img {
                width: 100%;
                height: 10rem;
                object-fit: contain;
                transition: transform 0.3s ease;
            }

            .gun-name {
                margin-top: 0.5rem;
                font-size: 0.875rem;
                color: #cccaca;
                text-align: center;
                width: 100%;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .selected-badge {
                position: absolute;
                top: 0.75rem;
                right: 0.75rem;
                width: 2rem;
                height: 2rem;
                background: #ad40ff;
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                box-shadow: 0 2px 10px rgba(173, 64, 255, 0.5);

                svg {
                    width: 1rem;
                    height: 1rem;
                    stroke: white;
                }
            }
        }

        .chromas {
            display: flex;
            flex-direction: row;
            gap: 0.5rem;
            padding: 0;
            flex-wrap: wrap;
            max-width: 20rem;
        }

        .chroma-item {
            position: relative;
            width: 4rem;
            height: 4rem;
            border-radius: 0.5rem;
            overflow: hidden;
            border: 2px solid rgba(60, 60, 60, 0.8);
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            background: rgba(0, 0, 0, 0.7);
            backdrop-filter: blur(10px);

            &:hover {
                border-color: #ad40ff;
                transform: scale(1.1);
                box-shadow: 0 4px 15px rgba(173, 64, 255, 0.5);
            }

            &.selected {
                border-color: #ad40ff;
                background: rgba(173, 64, 255, 0.3);
                box-shadow: 0 0 20px rgba(173, 64, 255, 0.6);
            }

            &:focus {
                outline: 2px solid #ad40ff;
                outline-offset: 2px;
            }

            .chroma-swatch {
                width: 100%;
                height: 100%;
                background-size: cover;
                background-position: center;
            }

            img {
                width: 100%;
                height: 100%;
                object-fit: contain;
            }

            .chroma-name {
                display: none;
            }

            .selected-badge {
                position: absolute;
                top: 0.25rem;
                right: 0.25rem;
                width: 1.25rem;
                height: 1.25rem;
                background: #ad40ff;
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                box-shadow: 0 2px 10px rgba(173, 64, 255, 0.5);

                svg {
                    width: 0.625rem;
                    height: 0.625rem;
                    stroke: white;
                }
            }
        }

        .levels {
            display: flex;
            flex-direction: row;
            gap: 0.5rem;
            padding: 0;
            flex-wrap: wrap;
            max-width: 20rem;
        }

        .level-item {
            position: relative;
            width: 3rem;
            height: 3rem;
            border-radius: 0.5rem;
            overflow: hidden;
            border: 2px solid rgba(60, 60, 60, 0.8);
            cursor: pointer;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(0, 0, 0, 0.7);
            backdrop-filter: blur(10px);

            &:hover {
                border-color: #ad40ff;
                transform: scale(1.1);
                box-shadow: 0 4px 15px rgba(173, 64, 255, 0.5);
            }

            &.selected {
                border-color: #ad40ff;
                background: rgba(173, 64, 255, 0.3);
                box-shadow: 0 0 20px rgba(173, 64, 255, 0.6);
            }

            &:focus {
                outline: 2px solid #ad40ff;
                outline-offset: 2px;
            }

            .level-icon {
                width: 100%;
                height: 100%;
                display: flex;
                align-items: center;
                justify-content: center;

                img {
                    width: 70%;
                    height: 70%;
                    object-fit: contain;
                }

                .level-number {
                    font-size: 0.875rem;
                    font-weight: 700;
                    color: #ffffff;
                }
            }

            .level-name {
                display: none;
            }

            .selected-badge {
                position: absolute;
                top: 0.125rem;
                right: 0.125rem;
                width: 1rem;
                height: 1rem;
                background: #ad40ff;
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                box-shadow: 0 2px 10px rgba(173, 64, 255, 0.5);

                svg {
                    width: 0.5rem;
                    height: 0.5rem;
                    stroke: white;
                }
            }
        }
    }
</style>
