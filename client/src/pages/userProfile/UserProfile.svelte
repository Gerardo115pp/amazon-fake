<script>
    import AddProductModal from './page-components/AddProductModal.svelte';
    import UserProfileHeader from './page-components/UserProfileHeader.svelte';
    import UserProducts from './page-components/UserProducts.svelte'
    import { server_name } from '../../server_info';
    import { onMount } from 'svelte';
    export let params = {};

    let { sk:session_key } = params;
    let is_porduct_modal_showing = false;
    let user_products = [];
    let user_data = {};

    onMount(async () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);
        
        //getting user data
        const request = new Request(`http://${server_name}/user`, {method: 'GET', headers: headers});
        user_data = await fetch(request).then(promise => promise.json());
        
        // getting user products
        requestUserProducts();
    });

    const requestUserProducts = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        fetch(new Request(`http://${server_name}/products`, {method: 'GET', headers: headers}))
            .then(promise => promise.json())
            .then(products => user_products = products);
    }

    const closeProductModal = () => is_porduct_modal_showing = !is_porduct_modal_showing

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
</style>

<div id="user-profile-page">
    {#if is_porduct_modal_showing}
        <AddProductModal {session_key} onProductAdded={requestUserProducts} close={closeProductModal}/>
    {/if}
    <UserProfileHeader {session_key} {user_data} {closeProductModal}/>
    <main id="user-profile-main">
        <aside id="profile-data">
            <UserProducts products={user_products}/>
        </aside>
        <div id="user-profile-panel">
            
        </div>
    </main>
</div>