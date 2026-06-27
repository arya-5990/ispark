<script lang="ts">
  import { fade } from 'svelte/transition';
  import { goto } from '$app/navigation';

  // State variables (Svelte 5 runes)
  let currentPassword = $state("");
  let newPassword = $state("");
  let confirmPassword = $state("");
  
  let showNewPassword = $state(false);
  let showConfirmPassword = $state(false);
  
  let submitting = $state(false);
  let success = $state(false);
  let errorMsg = $state("");

  // Deriving validation rules
  let isMinLength = $derived(newPassword.length >= 8);
  let hasUppercase = $derived(/[A-Z]/.test(newPassword));
  let hasLowercase = $derived(/[a-z]/.test(newPassword));
  let hasNumber = $derived(/\d/.test(newPassword));
  let hasSpecial = $derived(/[^A-Za-z0-9]/.test(newPassword));

  let isPasswordValid = $derived(isMinLength && hasUppercase && hasLowercase && hasNumber && hasSpecial);
  
  let isConfirmPasswordFilled = $derived(confirmPassword.length > 0);
  let isNewPasswordFilled = $derived(newPassword.length > 0);
  let doPasswordsMatch = $derived(
    isNewPasswordFilled &&
    isConfirmPasswordFilled &&
    newPassword === confirmPassword
  );

  let showMatchError = $derived(isNewPasswordFilled && isConfirmPasswordFilled && !doPasswordsMatch);

  // Rules list for rendering
  let passwordRules = $derived([
    { id: "minLength", label: "Minimum 8 Characters", met: isMinLength },
    { id: "uppercase", label: "One Uppercase Letter", met: hasUppercase },
    { id: "lowercase", label: "One Lowercase Letter", met: hasLowercase },
    { id: "number", label: "One Number", met: hasNumber },
    { id: "special", label: "One Special Character", met: hasSpecial },
  ]);

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    if (!isPasswordValid || !doPasswordsMatch || currentPassword.trim() === "") {
      return;
    }

    submitting = true;
    errorMsg = "";

    // Simulate password update API call
    setTimeout(() => {
      submitting = false;
      success = true;
    }, 1500);
  };

  const handleCancelLogout = () => {
    // Navigate back to the admin login page
    goto('/admin-portal');
  };
</script>

<div class="flex flex-col max-w-[560px] w-full items-center px-4 sm:px-6 lg:px-8 py-0 relative font-sans">
  {#if success}
    <!-- Success Card -->
    <div class="flex flex-col w-full max-w-[496px] items-center p-8 bg-white rounded-2xl border border-border-base shadow-sm text-center" in:fade={{ duration: 400 }}>
      <div class="w-16 h-16 bg-emerald-100 text-emerald-805 flex items-center justify-center rounded-full border-2 border-emerald-250 mb-6 animate-bounce">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor" class="w-8 h-8">
          <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
        </svg>
      </div>
      <h2 class="text-2xl font-bold text-slate-900 font-serif leading-tight">Password Updated!</h2>
      <p class="text-slate-550 text-sm mt-3 max-w-sm">
        Your credentials have been successfully updated. You can now use your new password to sign in.
      </p>
      <div class="mt-8 w-full">
        <a href="/admin-portal" class="flex w-full py-3.5 items-center justify-center bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-wider uppercase rounded-xl transition duration-200 shadow-xs">
          Return to Login
        </a>
      </div>
    </div>
  {:else}
    <!-- Security Policy Banner -->
    <section
      class="flex items-start gap-4 p-5 w-full max-w-[496px] bg-red-50/50 rounded-2xl border border-red-200/60 mb-6"
      aria-labelledby="security-policy-heading"
    >
      <div class="pt-0.5 shrink-0">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-5 h-5 text-acad-red">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
        </svg>
      </div>
      <div class="flex flex-col gap-1 select-none">
        <h2 id="security-policy-heading" class="text-sm font-bold text-acad-red uppercase tracking-wider">
          Security Policy
        </h2>
        <p class="text-xs text-red-800/80 leading-relaxed font-medium">
          Temporary passwords issued by administrators expire after first use. All mentor accounts must create a new password before continuing.
        </p>
      </div>
    </section>

    <!-- Main Setup Card -->
    <div class="flex flex-col w-full max-w-[496px] bg-white rounded-2xl border border-border-base shadow-sm overflow-hidden animate-fade-up">
      <!-- Card Header -->
      <div class="flex flex-col items-start px-6 sm:px-8 py-6 border-b border-border-base bg-slate-50/50">
        <div class="px-2.5 py-0.5 bg-slate-200/70 rounded text-[10px] font-bold tracking-widest text-inst-navy uppercase">
          PASSWORD SETUP
        </div>
        <h1 class="text-2xl font-bold text-inst-navy font-serif leading-tight mt-2">
          Create Secure Password
        </h1>
        <p class="text-slate-500 text-xs mt-1">
          Update your credentials to continue to the Mentor Portal
        </p>
      </div>

      <!-- Form -->
      <form onsubmit={handleSubmit} class="p-6 sm:p-8 flex flex-col gap-5">
        
        <!-- Current Password -->
        <div class="flex flex-col gap-1.5">
          <label for="currentPassword" class="text-[11px] font-bold text-slate-700 tracking-wider uppercase">
            Current Password
          </label>
          <div class="relative flex items-center">
            <div class="absolute left-3.5 pointer-events-none">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-slate-400">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 0 1 3 3m3 0a6 6 0 0 1-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1 1 21.75 8.25Z" />
              </svg>
            </div>
            <input
              id="currentPassword"
              name="currentPassword"
              type="password"
              autoComplete="current-password"
              bind:value={currentPassword}
              placeholder="Enter temporary password"
              class="w-full pl-10 pr-4 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
              required
            />
          </div>
        </div>

        <!-- New Password -->
        <div class="flex flex-col gap-1.5">
          <label for="newPassword" class="text-[11px] font-bold text-slate-700 tracking-wider uppercase">
            New Password
          </label>
          <div class="relative flex items-center">
            <div class="absolute left-3.5 pointer-events-none">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-slate-400">
                <rect width="18" height="11" x="3" y="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
            </div>
            <input
              id="newPassword"
              name="newPassword"
              type={showNewPassword ? "text" : "password"}
              autoComplete="new-password"
              bind:value={newPassword}
              placeholder="Create a new password"
              class="w-full pl-10 pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
              aria-describedby="password-requirements"
              required
            />
            <button
              type="button"
              onclick={() => showNewPassword = !showNewPassword}
              class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-650 focus:outline-none"
              aria-label={showNewPassword ? "Hide new password" : "Show new password"}
            >
              {#if showNewPassword}
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" />
                </svg>
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                </svg>
              {/if}
            </button>
          </div>
        </div>

        <!-- Confirm Password -->
        <div class="flex flex-col gap-1.5">
          <label for="confirmPassword" class="text-[11px] font-bold text-slate-700 tracking-wider uppercase">
            Confirm Password
          </label>
          <div class="relative flex items-center">
            <div class="absolute left-3.5 pointer-events-none">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 text-slate-400">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" />
              </svg>
            </div>
            <input
              id="confirmPassword"
              name="confirmPassword"
              type={showConfirmPassword ? "text" : "password"}
              autoComplete="new-password"
              bind:value={confirmPassword}
              placeholder="Re-enter new password"
              class="w-full pl-10 pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
              required
            />
            <button
              type="button"
              onclick={() => showConfirmPassword = !showConfirmPassword}
              class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-650 focus:outline-none"
              aria-label={showConfirmPassword ? "Hide confirm password" : "Show confirm password"}
            >
              {#if showConfirmPassword}
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" />
                </svg>
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                </svg>
              {/if}
            </button>
          </div>
        </div>

        <!-- Password Requirements Checklist -->
        <div
          id="password-requirements"
          class="flex flex-col gap-2.5 p-4 bg-slate-50/50 rounded-xl border border-border-base/70"
          aria-live="polite"
        >
          {#each passwordRules as rule}
            <div class="flex items-center gap-2.5 text-xs">
              <div class="flex shrink-0">
                {#if rule.met}
                  <!-- Met: Green Check -->
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="3" stroke="currentColor" class="w-3.5 h-3.5 text-emerald-600">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                  </svg>
                {:else}
                  <!-- Not Met: Gray Circle -->
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-3.5 h-3.5 text-slate-300">
                    <circle cx="12" cy="12" r="9" />
                  </svg>
                {/if}
              </div>
              <span class={rule.met ? "text-slate-800 font-medium" : "text-slate-450"}>
                {rule.label}
              </span>
            </div>
          {/each}
        </div>

        <!-- Error Message / Passwords Match Warning -->
        {#if showMatchError}
          <div class="flex items-start gap-2 text-xs text-acad-red/90" in:fade={{ duration: 150 }}>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4 mt-0.5 shrink-0">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
            </svg>
            <p>Confirm password must match the new password.</p>
          </div>
        {/if}

        <!-- Actions -->
        <div class="flex flex-col gap-4 mt-2">
          <button
            type="submit"
            disabled={submitting || !isPasswordValid || !doPasswordsMatch}
            class="w-full py-3.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-wider uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
          >
            {#if submitting}
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Updating...
            {:else}
              Update Password
            {/if}
          </button>

          <button
            type="button"
            onclick={handleCancelLogout}
            class="text-xs font-bold text-slate-500 hover:text-slate-800 transition-colors focus:outline-none py-1.5 self-center"
          >
            Cancel & Logout
          </button>
        </div>
      </form>
    </div>
  {/if}
</div>

<style>
  .animate-fade-up {
    opacity: 0;
    transform: translateY(25px);
    animation: fadeUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  }

  @keyframes fadeUp {
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
