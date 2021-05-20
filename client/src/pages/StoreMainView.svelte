<script>
    import StoreTopBar from '../components/StoreTopBar.svelte';

    export let params = {};
    let { sk:session_key } = params;
</script>

<div id="store-entry-proint">
    <StoreTopBar {session_key}/>
</div>