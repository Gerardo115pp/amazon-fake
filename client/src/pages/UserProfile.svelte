<script>
    import AddProductModal from '../components/AddProductModal.svelte';
    import goback_svg from '../resources/return.svg';
    import { server_name } from '../server_info';
    import { pop } from 'svelte-spa-router';
    import add_svg from '../resources/add.svg';
    import { onMount } from 'svelte';
    export let params = {};

    let { sk:session_key } = params;
    let is_porduct_modal_showing = false;
    let user_data = {};

    onMount(async () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);
        const request = new Request(`http://${server_name}/user`, {method: 'GET', headers: headers});
        user_data = await fetch(request).then(promise => promise.json());
    });

    const closeProductModal = () => is_porduct_modal_showing = !is_porduct_modal_showing

</script>

<style>
    #profile-controls {
        display: flex;
        width: 97vw;
        padding: 0.5% 1%;
        justify-content: space-between;
        align-items: center;
    }

    :global(#profile-controls svg) {
        cursor: pointer;
        width: 3vw;
        fill: var(--theme-color);
    }

    #user-name-label {
        font-size: 1.3rem;
    }

    #add-product {
        margin-left: 80%;
    }

    :global(#add-product svg) {
        width: 2vw;
    }
</style>

<div id="user-profile-page">
    {#if is_porduct_modal_showing}
        <AddProductModal close={closeProductModal}/>
    {/if}
    <header id="profile-controls">
        <div on:click={() => pop()} class="goback-container">
            {@html goback_svg}
        </div>
        <div on:click={closeProductModal} id="add-product">
            {@html add_svg}
        </div>
        <div id="user-name-label">
            <span>{user_data.username}</span>
        </div>
    </header>
</div>