<script>
    import StoreTopBar from './page-components/StoreTopBar.svelte';
    import ProductFeed from './page-components/ProductsFeed.svelte';
    import { onMount } from 'svelte';
    import { server_name } from '../../server_info';

    let feed_products = [];

    export let params = {};
    let { sk:session_key } = params;

    onMount(() => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/products?id=*`, { method:'GET', headers: headers});
        fetch(request)
            .then(promise => promise.json())
            .then(products => feed_products = products);
    });

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
    <StoreTopBar {session_key}/>
    <main id="store-main-content">
        <ProductFeed products={feed_products}/>
    </main>
</div>