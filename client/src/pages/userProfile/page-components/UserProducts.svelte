<script>
    import FieldData  from '../../../classes/FieldData';
import { user_events } from '../../../events';
    import { server_name } from '../../../server_info';

    export let products = [];
    export let session_key = "";

    let restock_field = new FieldData("restock-product", /\d+/, 'restock');
    let selected_item = -1;

    const awaitForEnter = e => {
        if(e.key.toLowerCase() === "enter") {
            restock();
        }
    }

    const restock = () => {
        if (restock_field.isReady()) {
            const headers = new Headers();
            headers.set('X-sk', session_key);

            const form_data = new FormData();
            form_data.append("restock_by", restock_field.getFieldValue());
            form_data.append("product_id", products[selected_item].id);

            const request = new Request(`http://${server_name}/products`, {method: 'PATCH', body: form_data, headers: headers});
            fetch(request)
                .then(promise => {
                    if (promise.ok) {
                        window.dispatchEvent(new Event(user_events.PRODUCTS_CHANGED));
                        selected_item = -1
                    }
                })
        }
    }

</script>

<style>
    #user-products-container {
        background: var(--theme-color);
        width: 20vw;
        height: 60vh;
        color: white;
        padding: 1vh 1vw;
        border-radius: 15px;
    }

    .user-product-container {
        cursor: pointer;
        display: flex;
        min-height: 5vh;
        max-height: 10vh;
        justify-content: space-around;
        flex-wrap: wrap;
        border-bottom: 1px solid white;
        transition: background .2s ease-in-out;
    }

    .user-product-container:hover {
        background-color: rgba(255, 255, 255, 0.5);
    }

    .user-product-container * {
        user-select: none;
    }

    #resuplier-container {
        width: 20vw;
        padding: 0 1vw .5vh;
    }

    #resuplier-container input{
        font-family: var(--main-font);
        color: var(--theme-color);
        border: none;
        border-radius: 5px;
    }
</style>

<div id="user-products-container">
    {#each products as p,h}
        <div on:click={() => selected_item = h} class="user-product-container">
            <span class="product-name">{p.name}</span>
            <span class="product-price">{p.price.toLocaleString('en')} MXN</span>
            <span class="product-sales">{p.solds}/{p.stock}</span>
            {#if h === selected_item}
                <div id="resuplier-container">
                    <input on:keydown={awaitForEnter} id="{restock_field.id}" placeholder={restock_field.name} type="text">
                </div>
            {/if}
        </div>
    {/each}
</div>