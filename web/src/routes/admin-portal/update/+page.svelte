<script lang="ts">
  import HeaderSection from '../HeaderSection.svelte';
  import UpdatePasswordSection from './UpdatePasswordSection.svelte';
  import FooterSection from '../FooterSection.svelte';
</script>

<svelte:head>
  <title>Setup Admin Password | iSPARC IIPS DAVV</title>
  <meta name="description" content="Update your administrator credentials to access the iSPARC administration panel." />
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="min-h-screen bg-bg-main flex flex-col items-center w-full">
  <HeaderSection />
  <main class="flex-grow w-full flex justify-center items-start py-12">
    <UpdatePasswordSection />
  </main>
  <FooterSection />
</div>
