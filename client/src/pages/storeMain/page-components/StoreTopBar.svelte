<script>
    
    import user_profile_logo from '../../../resources/man.svg';
    import cart_svg from '../../../resources/cart.svg';
    import { push } from 'svelte-spa-router';

    export let session_key = "";
</script>

<style>
    #control-navbar {
        display: flex;
        background: var(--theme-gradiant);
        height: 7vh;
        justify-content: space-between;
        align-items: center;
    }

    #cn-searchbar {
        display: flex;
        margin-left: 15vw;
    }

    #cn-searchbar input {
        font-family: var(--main-font);
        width: 30vw;
        color: var(--theme-color);
        font-size: 1.1rem;
        padding: 0.2rem 0.4rem;
        border: none;
        border-radius: 3px;
        outline: none;
    }

    #cn-controls {
        display: flex;
    }

    .link-container {
        width: 2.3vw;
        margin-right: 5vw;
    }

    .link-container span{
        cursor: pointer;
        display: flex;
        align-items: center;
    }

    :global(.link-container svg) { 
        fill: white;
    }
</style>

<div id="control-navbar">
    <div id="cn-searchbar" class="cn-item">
        <input placeholder="search" type="search">
    </div>
    <div id="cn-controls" class="controls cn-item">
        <div on:click={() => push(`/cart/${session_key}`)} class="link-container" id="cart-link">
            <span id="cart-logo">
                {@html cart_svg}
            </span>
        </div>
        <div on:click={() => push(`/profile/${session_key}`)} class="link-container" id="profile-link">
            <span id="profile-logo">
                {@html user_profile_logo}
            </span>
        </div>
    </div>
</div>