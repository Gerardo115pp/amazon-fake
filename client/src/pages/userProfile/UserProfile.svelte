<script>
    import AddProductModal from './page-components/AddProductModal.svelte';
    import UserProfileHeader from './page-components/UserProfileHeader.svelte';
    import UserProducts from './page-components/UserProducts.svelte';
    import { user_events } from '../../events';
    import { server_name } from '../../server_info';
    import { onMount } from 'svelte';
    export let params = {};

    let { sk:session_key } = params;
    let is_porduct_modal_showing = false;
    let user_products = [];
    let user_data = {};
    let user_stash = 0;

    onMount(async () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);
        
        //getting user data
        const request = new Request(`http://${server_name}/user`, {method: 'GET', headers: headers});
        user_data = await fetch(request).then(promise => promise.json());
        
        // getting user products
        requestUserProducts();

        // getting user stash
        requestUserStash();

        // setting event lisenters
        setEventLisenters();
    });

    const closeProductModal = () => is_porduct_modal_showing = !is_porduct_modal_showing

    const requestUserStash = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/user-stash`, { method: 'GET', headers: headers});
        fetch(request)
            .then(promise => promise.json())
            .then(response => {
                if (response.response) {
                    user_stash = response.response;
                }
            });
    }

    const requestUserProducts = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        fetch(new Request(`http://${server_name}/products`, {method: 'GET', headers: headers}))
            .then(promise => promise.json())
            .then(products => user_products = products);
    }

    const setEventLisenters = () => {
        window.addEventListener(user_events.PRODUCTS_CHANGED, requestUserProducts);
    }    

</script>

<style>
    #user-profile-main {
        display: flex;
        height: 100vh;
    }

    #profile-data {
        display: flex;
        width: 30vw;
        margin-left: auto;
        justify-content: center;
        align-items: center;
    }

    #user-profile-panel {
        width: 70vw;
    }

    #user-profile-status-bar {
        display: flex;
        height: 10vh;
        font-size: 1.2rem;
        font-weight: bolder;
        padding: 0 4vw;
        align-items: center;
    }
</style>

<div id="user-profile-page">
    {#if is_porduct_modal_showing}
        <AddProductModal {session_key} onProductAdded={requestUserProducts} close={closeProductModal}/>
    {/if}
    <UserProfileHeader {session_key} {user_data} {closeProductModal}/>
    <main id="user-profile-main">
        <div id="user-profile-panel">
            <div id="user-profile-status-bar">
                <span id="user-stash">
                    Earning: ${user_stash.toLocaleString('en')} MXN
                </span>
            </div>
        </div>
        <aside id="profile-data">
            <UserProducts {session_key} products={user_products}/>
        </aside>
    </main>
</div>