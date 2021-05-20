<script>
    export let label = "redirect";
    export let onClick = trash => {};
    label = label.toLowerCase();
</script>

<style>
    .option-redirect {
        width: 20%;
        background: var(--theme-gradiant);
        color: var(--dark-color);
        text-align: center;
        text-transform: capitalize;
        border-radius: 9999px;
        padding: .5em;
        transition: all .3s ease-in-out;
    }
    
    .option-redirect span{
        display: block;
        user-select: none;
        transition: all .2s ease-in-out;
    }

    .option-redirect:hover {
        filter: brightness(1.2);
    }

    .option-redirect:hover span {
        transform: scale(1.3);
        color: antiquewhite;
    }

</style>

<div on:click={() => onClick(label)} class="option-redirect">
    <span>{label}</span>
</div>