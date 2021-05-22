<script>
    import CartProduct from './page-components/CartProduct.svelte'
    import CartHero from './page-components/CartHero.svelte';
    import Checkout from './page-components/Checkout.svelte';
    import { server_name } from '../../server_info';
    import { pop } from 'svelte-spa-router';
    import { onMount } from 'svelte';

    export let params = {};
    
    let { sk:session_key } = params;
    let cart_content = [];

    onMount(() => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/cart`, {method: 'GET', headers: headers});

        fetch(request)
            .then(promise => promise.json())
            .then(cart => cart_content = cart);

    });
</script>

<style>
    #cart-page-container {
        display: flex;
        width: 100%;
        height: 100vh;
    }

    #cart-content {
        width: 70%;
    }

    #hero-container {
        display: flex;
        height: 20vh;
        align-items: center;
    }

    #cart-products-container {
        display: flex;
        height: 70vh;
        justify-content: center;
        align-items: center;
    }

    #cart-products {
        width: 50vw;
        height: 50vh;
    }

    #nav-footer {
        display: flex;
        height: 6vh;
        padding: 2vh 2vw;
        align-items: center;
    }

    .fn-link {
        cursor: pointer;
    }

    #checkout-section {
        width: 30%;
        height: 100vh;
        background: var(--theme-gradiant);
        border-radius: 199px 0 0 0px;
        padding-top: 20vh;
        box-shadow: 0 0 20px 25px var(--modal-shadows);
    }

</style>

<div id="cart-page-container">
    <div id="cart-content">
        <div id="hero-container">
            <CartHero/>
        </div>
        <div id="cart-products-container">
            <div id="cart-products">
                {#each cart_content as c}
                    <CartProduct {session_key} updateCart={new_cart => cart_content = new_cart} product_data={c}/>
                {/each}
            </div>
        </div>
        <footer id="nav-footer">
            <div on:click={() => pop()} class="fn-link">Go back</div>
        </footer>
    </div>
    <div id="checkout-section">
        <Checkout {session_key} products_data={cart_content}/>
    </div>
</div>