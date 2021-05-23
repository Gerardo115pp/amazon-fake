<script>
    import StoreTopBar from './page-components/StoreTopBar.svelte';
    import PurcheasModal from './page-components/StoreBuyProductModal.svelte'
    import ProductFeed from './page-components/ProductsFeed.svelte';
    import Transmisor from '../../classes/Trasmisor';
    import { server_name } from '../../server_info';
    import { onMount, onDestroy } from 'svelte';

    let selected_product_data = {};
    let is_purchease_modal_showing = false;
    let feed_products = [];
    let transmisor = new Transmisor();

    export let params = {};
    let { sk:session_key } = params;

    onMount(() => {
        updateProducts();

        // setting transmisor
        transmisor.onMessage = () => updateProducts();
        transmisor.connect(`ws://${server_name}/products-feed?sk=${session_key}`);
    });

    onDestroy(() => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/products-feed`, {method: 'DELETE', headers: headers});
        fetch(request)
    })

    const updateProducts = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/products?id=*`, { method:'GET', headers: headers});
        fetch(request)
            .then(promise => promise.json())
            .then(products => {
                feed_products = products
            });
    }

    const openPurchaseModal = product_data => {
        selected_product_data = product_data;
        is_purchease_modal_showing = true;
    }
</script>

<style>
    #store-main-content {
        display: flex;
        height: 93vh;
        justify-content: center;
        align-items: center;
    }
</style>

<div id="store-entry-proint">
    {#if is_purchease_modal_showing}
        <PurcheasModal product_data={selected_product_data} close={() => is_purchease_modal_showing = false} {session_key}/>
    {/if}
    <StoreTopBar {session_key}/>
    <main id="store-main-content">
        <ProductFeed openProduct={openPurchaseModal} products={feed_products}/>
    </main>
</div>