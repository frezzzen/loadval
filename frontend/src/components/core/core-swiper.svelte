<script lang="ts">
    import Swiper from "swiper/bundle";
    import "swiper/css";
    import "swiper/css/bundle";
    import { onDestroy, onMount, type Snippet, tick } from "svelte";
    import type { SwiperOptions } from "swiper/types";

    type Props = {
        children: Snippet;
        options?: SwiperOptions;
        activeElement?: HTMLElement | undefined;
        activeElementIndex?: number | undefined;
        defaultNavArrows?: boolean;
        defaultPagination?: boolean;
        arrowLeftIcon?: Snippet;
        arrowRightIcon?: Snippet;
        onSlideChange?: (index: number) => void;
    };

    let {
        children,
        options,
        defaultNavArrows = false,
        defaultPagination = false,
        activeElement = $bindable(),
        activeElementIndex = $bindable(),
        arrowLeftIcon,
        arrowRightIcon,
        onSlideChange,
    }: Props = $props();

    let swiper: Swiper;

    function createSwiper() {
        if (swiper) swiper.destroy(true, true);

        swiper = new Swiper("#swiper", {
            slidesPerView: 1,
            spaceBetween: 0,
            initialSlide: activeElementIndex ?? 0,
            pagination: {
                el: ".swiper-pagination",
                clickable: true,
            },
            autoplay: {
                delay: 2500,
                disableOnInteraction: false,
            },
            navigation: {
                nextEl: ".swiper-button-next",
                prevEl: ".swiper-button-prev",
            },
            observer: true,
            observeParents: true,
            observeSlideChildren: true,
            ...options,
        });

        // swiper.on("init", () => {
        //   activeElement = swiper.slides.at(-1)
        // })

        swiper.on("slideChange", function () {
            activeElement = swiper.slides.at(swiper.activeIndex - 1);
            activeElementIndex = swiper.activeIndex;
            onSlideChange?.(swiper.activeIndex);
        });

        swiper.init();
    }

    onMount(async () => {
        createSwiper();
        window.addEventListener("resize", () => {
            swiper?.update();
            const isAutoplayEnabled = (
                swiper?.params?.autoplay as { enabled: boolean } | undefined
            )?.enabled;
            if (isAutoplayEnabled) {
                swiper?.autoplay?.start();
            } else {
                swiper?.autoplay?.stop();
            }
        });
        await tick();
        if (activeElementIndex === undefined) {
            activeElementIndex = 0;
        } else {
            swiper.slideTo(activeElementIndex);
        }
    });

    onDestroy(() => {
        if (swiper) swiper.destroy(true, true);
        window.removeEventListener("resize", swiper.update);
    });
</script>

<div id="swiper" class="swiper">
    <div class="swiper-wrapper">
        {@render children()}
    </div>

    <div
        class="swiper-pagination"
        style="display: {defaultPagination ? 'flex' : 'none'}"
    ></div>
    <div
        class="swiper-button-prev"
        style="display: {defaultNavArrows ? 'flex' : 'none'}"
    >
        {#if arrowLeftIcon}
            {@render arrowLeftIcon()}
        {/if}
    </div>
    <div
        class="swiper-button-next"
        style="display: {defaultNavArrows ? 'flex' : 'none'}"
    >
        {#if arrowRightIcon}
            {@render arrowRightIcon()}
        {/if}
    </div>
</div>

<style lang="scss">
    .swiper {
        position: relative;
        width: 100%;
        height: 100%;
        overflow: hidden;
    }
</style>
