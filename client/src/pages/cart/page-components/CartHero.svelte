<script>
    import cart_logo from '../../../resources/checkout.svg';
</script>

<style>
    #cart-hero-container {
        display: flex;
        width: 50%;
        height: 10vh;
        justify-content: space-evenly;
    }

    #logo-container {
        display: flex;
        width: 10%;
        align-items: center;
    }

    :global(#logo-container svg) {
        fill: var(--ready);
    }

    #separator {
        width: 1%;
        min-height: 100%;
        background: var(--theme-color);
    }

    #title-container {
        font-size: 1.8rem;
        user-select: none;
        margin: auto 0;
    }

</style>

<div id="cart-hero-container">
    <div id="logo-container">
        {@html cart_logo}
    </div>
    <div id="separator"></div>
    <div id="title-container">Your shopping cart</div>
</div>