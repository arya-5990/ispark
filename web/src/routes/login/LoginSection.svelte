<script lang="ts">
	import { fade, slide } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { API_BASE_URL } from '$lib/config';

	const formId = 'student-login';

	// Reactive state for the login form fields
	let email = $state('');
	let password = $state('');
	let rememberMe = $state(false);
	let showPassword = $state(false);

	// API and UI states
	let submitting = $state(false);
	let loginSuccess = $state(false);
	let errorMsg = $state('');

	// Captcha states
	let captchaQuestion = $state('');
	let captchaID = $state('');
	let captchaAnswer = $state('');
	let captchaLoading = $state(false);

	// Forgot Password / Reset Password states
	let viewState = $state<'login' | 'forgot' | 'reset'>('login');
	let forgotEmail = $state('');
	let forgotOtp = $state('');
	let forgotNewPassword = $state('');
	let forgotConfirmPassword = $state('');
	let showForgotNewPassword = $state(false);
	let showForgotConfirmPassword = $state(false);

	// Form validity derived state
	let isFormValid = $derived(
		email.trim() !== '' && password.trim() !== '' && captchaAnswer.trim() !== ''
	);

	// Password Validation Runes (for reset password)
	let resetHasMinLength = $derived(forgotNewPassword.length >= 8);
	let resetHasUppercase = $derived(/[A-Z]/.test(forgotNewPassword));
	let resetHasLowercase = $derived(/[a-z]/.test(forgotNewPassword));
	let resetHasNumber = $derived(/[0-9]/.test(forgotNewPassword));
	let resetHasSpecialChar = $derived(/[^A-Za-z0-9]/.test(forgotNewPassword));

	let resetAllPasswordRequirementsMet = $derived(
		resetHasMinLength &&
			resetHasUppercase &&
			resetHasLowercase &&
			resetHasNumber &&
			resetHasSpecialChar
	);

	let resetPasswordsMatch = $derived(
		forgotNewPassword !== '' && forgotNewPassword === forgotConfirmPassword
	);

	const featureItems = [
		{ text: 'Browse approved activities' },
		{ text: 'Manage course enrollments' },
		{ text: 'Upload and verify certificates' },
		{ text: 'Track credit progress' },
		{ text: 'Download extracurricular marksheet' },
		{ text: 'View leaderboard rankings' }
	];

	const importantNoticeItems = [
		'Login is permitted only for registered students.',
		'Use your institutional email address.',
		'Multiple failed login attempts may temporarily lock your account.',
		'Contact the iSPARC Coordinator if you cannot access your account.'
	];

	const helpItems = [
		'Forgot your password?',
		"Didn't receive verification email?",
		'Having trouble logging in?'
	];

	// Fetch Captcha
	async function fetchCaptcha() {
		captchaLoading = true;
		try {
			const response = await fetch(`${API_BASE_URL}/api/auth/captcha`);
			if (response.ok) {
				const data = await response.json();
				captchaID = data.captcha_id;
				captchaQuestion = data.question;
				captchaAnswer = '';
			}
		} catch (err) {
			console.error('Error fetching captcha:', err);
		} finally {
			captchaLoading = false;
		}
	}

	// Load captcha on mount
	$effect(() => {
		fetchCaptcha();
	});

	// Submit Handler
	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (!isFormValid) return;

		submitting = true;
		errorMsg = '';

		try {
			const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email_id: email,
					password: password,
					captcha_id: captchaID,
					captcha_answer: captchaAnswer
				})
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Invalid credentials');
			}

			if (data.access_token) {
				localStorage.setItem('access_token', data.access_token);
			}
			loginSuccess = true;
			setTimeout(() => {
				goto('/portal');
			}, 1500);
		} catch (err) {
			errorMsg =
				err instanceof Error ? err.message : 'Failed to login. Please check your credentials.';
			// Refresh captcha on login failure
			fetchCaptcha();
		} finally {
			submitting = false;
		}
	}

	// Handle Forgot Password Request
	async function handleForgotPassword(event: SubmitEvent) {
		event.preventDefault();
		if (forgotEmail.trim() === '') return;

		submitting = true;
		errorMsg = '';

		try {
			const response = await fetch(`${API_BASE_URL}/api/auth/forgot-password`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email_id: forgotEmail
				})
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Failed to request password reset');
			}

			viewState = 'reset';
		} catch (err) {
			errorMsg =
				err instanceof Error ? err.message : 'Failed to process request. Please try again.';
		} finally {
			submitting = false;
		}
	}

	// Handle Reset Password Submission
	async function handleResetPassword(event: SubmitEvent) {
		event.preventDefault();
		if (!resetAllPasswordRequirementsMet) {
			errorMsg = 'Password must meet all requirements.';
			return;
		}
		if (!resetPasswordsMatch) {
			errorMsg = 'Passwords do not match.';
			return;
		}

		submitting = true;
		errorMsg = '';

		try {
			const response = await fetch(`${API_BASE_URL}/api/auth/reset-password`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email_id: forgotEmail,
					otp_code: forgotOtp.trim(),
					password: forgotNewPassword,
					confirm_password: forgotConfirmPassword
				})
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Failed to reset password');
			}

			// Success
			viewState = 'login';
			errorMsg = '';
			alert('Password reset successfully! Please login with your new password.');
		} catch (err) {
			errorMsg = err instanceof Error ? err.message : 'Failed to reset password. Please try again.';
		} finally {
			submitting = false;
		}
	}
</script>

<section
	class="max-w-6xl w-full mx-auto grid grid-cols-1 lg:grid-cols-12 gap-10 px-4 sm:px-6 lg:px-8 py-12 items-start font-sans"
>
	<!-- Left Column: Guidelines & Information -->
	<aside class="lg:col-span-5 flex flex-col gap-6 animate-fade-in">
		<div class="flex flex-col gap-4">
			<div>
				<h1 class="text-4xl font-extrabold text-acad-red font-serif leading-tight tracking-tight">
					Login to your<br />
					<span class="text-inst-navy">iSPARC Account</span>
				</h1>
			</div>

			<p class="text-slate-500 text-sm leading-relaxed">
				Sign in using your registered institutional credentials to access your extracurricular
				records, enrolled activities, certificates, credits, leaderboard ranking, and academic
				progress.
			</p>
		</div>

		<!-- Feature list -->
		<div class="flex flex-col gap-3.5 bg-white p-6 rounded-xl border border-border-base shadow-xs">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">
				Student Portal Features
			</h3>
			{#each featureItems as item}
				<div class="flex items-start gap-3 group">
					<div
						class="flex w-5 h-5 items-center justify-center shrink-0 bg-emerald-50 rounded-full border border-emerald-250 transition-colors group-hover:bg-emerald-100"
					>
						<!-- Green checkmark SVG -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="3"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="w-3.5 h-3.5 text-emerald-600"
						>
							<polyline points="20 6 9 17 4 12" />
						</svg>
					</div>
					<span
						class="text-[13px] text-slate-700 font-medium group-hover:text-slate-900 transition-colors"
					>
						{item.text}
					</span>
				</div>
			{/each}
		</div>

		<!-- Security notice -->
		<div class="flex items-start gap-3.5 p-4 bg-amber-50/50 rounded-xl border border-amber-200/60">
			<div class="shrink-0 p-1 bg-amber-100/80 rounded-lg text-amber-800">
				<!-- Shield-alert SVG -->
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="w-5 h-5"
				>
					<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
					<line x1="12" y1="8" x2="12" y2="12" />
					<line x1="12" y1="16" x2="12.01" y2="16" />
				</svg>
			</div>
			<div>
				<h4 class="text-xs font-semibold text-slate-800 leading-tight">Security Notice</h4>
				<p class="text-[11px] text-slate-600 mt-1 leading-relaxed">
					Only registered students with verified institutional accounts may access the iSPARC
					portal. All login sessions are protected using secure authentication.
				</p>
			</div>
		</div>

		<!-- Important Notice -->
		<section
			class="flex flex-col gap-3 p-5 bg-white rounded-xl border border-border-base shadow-xs"
		>
			<h2 class="text-xs font-bold text-slate-400 tracking-widest uppercase">IMPORTANT NOTICE</h2>
			<ul class="flex flex-col gap-2.5 list-none">
				{#each importantNoticeItems as item}
					<li class="flex items-start gap-2.5 text-xs text-slate-600">
						<span class="text-acad-red font-bold text-sm leading-none">•</span>
						<p class="leading-relaxed">{item}</p>
					</li>
				{/each}
			</ul>
		</section>
	</aside>

	<!-- Right Column: Login Card -->
	<main
		class="lg:col-span-7 bg-white rounded-2xl border border-border-base shadow-sm overflow-hidden animate-fade-up"
	>
		{#if loginSuccess}
			<!-- Success View -->
			<div
				class="p-8 sm:p-12 text-center flex flex-col items-center gap-6"
				in:fade={{ duration: 400 }}
			>
				<div
					class="w-16 h-16 bg-emerald-100 text-emerald-800 flex items-center justify-center rounded-full border-2 border-emerald-250 animate-bounce"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="3"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="w-8 h-8"
					>
						<polyline points="20 6 9 17 4 12" />
					</svg>
				</div>

				<div>
					<h2 class="text-2xl font-bold text-slate-900 font-serif">Logged In Successfully!</h2>
					<p class="text-slate-500 text-sm mt-2 max-w-md mx-auto">
						Welcome back! You have successfully authenticated as <span
							class="font-semibold text-slate-800">{email}</span
						>. Redirecting you to your student portal dashboard...
					</p>
				</div>

				<div class="flex gap-3 w-full max-w-sm pt-4">
					<a
						href="/portal"
						class="flex-1 py-3 text-center bg-inst-navy hover:bg-inst-navy/90 text-white font-semibold text-xs tracking-wider uppercase rounded-xl transition duration-200"
					>
						Go to Dashboard
					</a>
				</div>
			</div>
		{:else if viewState === 'forgot'}
			<!-- Forgot Password View -->
			<div class="p-6 sm:p-8 border-b border-border-base bg-slate-50/50">
				<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
					STUDENT PORTAL
				</div>
				<h2 class="text-2xl font-bold text-inst-navy font-serif leading-tight mt-1">
					Reset Password
				</h2>
				<p class="text-slate-500 text-xs mt-1">
					Enter your registered institutional email to receive a reset code.
				</p>
			</div>

			<form onsubmit={handleForgotPassword} class="p-6 sm:p-8 flex flex-col gap-6">
				{#if errorMsg}
					<div
						class="p-3.5 bg-rose-50 border border-rose-200 text-rose-700 text-xs font-semibold rounded-lg"
						transition:slide={{ duration: 200 }}
					>
						{errorMsg}
					</div>
				{/if}

				<!-- Email -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-forgot-email"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						EMAIL
					</label>
					<div class="relative flex items-center">
						<span class="absolute left-3.5 pointer-events-none text-slate-400">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								style="margin-left: 10px"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2.2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75"
								/>
							</svg>
						</span>
						<input
							id="{formId}-forgot-email"
							type="email"
							bind:value={forgotEmail}
							placeholder="student.iips@gmail.com"
							class="w-full pr-3.5 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							style="padding-left: 40px;"
							required
						/>
					</div>
				</div>

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={submitting || forgotEmail.trim() === ''}
					class="w-full py-3.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-widest uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
				>
					{#if submitting}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						SENDING CODE...
					{:else}
						SEND RESET CODE
					{/if}
				</button>

				<div class="flex justify-center items-center text-xs mt-2">
					<button
						type="button"
						onclick={() => {
							viewState = 'login';
							errorMsg = '';
						}}
						class="text-slate-500 hover:text-slate-800 font-semibold transition duration-200"
					>
						← Back to Login
					</button>
				</div>
			</form>
		{:else if viewState === 'reset'}
			<!-- Reset Password View -->
			<div class="p-6 sm:p-8 border-b border-border-base bg-slate-50/50">
				<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
					STUDENT PORTAL
				</div>
				<h2 class="text-2xl font-bold text-inst-navy font-serif leading-tight mt-1">
					Enter Reset Details
				</h2>
				<p class="text-slate-500 text-xs mt-1">
					Enter the verification code sent to {forgotEmail} and choose a new password.
				</p>
			</div>

			<form onsubmit={handleResetPassword} class="p-6 sm:p-8 flex flex-col gap-6">
				{#if errorMsg}
					<div
						class="p-3.5 bg-rose-50 border border-rose-200 text-rose-700 text-xs font-semibold rounded-lg"
						transition:slide={{ duration: 200 }}
					>
						{errorMsg}
					</div>
				{/if}

				<!-- OTP Code -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-reset-otp"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						VERIFICATION CODE (OTP)
					</label>
					<input
						id="{formId}-reset-otp"
						type="text"
						bind:value={forgotOtp}
						placeholder="Enter Code"
						class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-border-base text-center tracking-widest text-[15px] font-bold text-slate-800 placeholder:text-slate-400 placeholder:tracking-normal focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
						required
					/>
				</div>

				<!-- New Password -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-reset-new-password"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						NEW PASSWORD
					</label>
					<div class="relative flex items-center">
						<span class="absolute left-3.5 pointer-events-none text-slate-400">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								style="margin-left: 10px"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2.2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
								/>
							</svg>
						</span>
						<input
							id="{formId}-reset-new-password"
							type={showForgotNewPassword ? 'text' : 'password'}
							bind:value={forgotNewPassword}
							placeholder="Enter new password"
							class="w-full pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							style="padding-left: 40px;"
							required
						/>
						<button
							type="button"
							onclick={() => (showForgotNewPassword = !showForgotNewPassword)}
							class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none"
							aria-label={showForgotNewPassword ? 'Hide password' : 'Show password'}
						>
							{#if showForgotNewPassword}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<line x1="2" y1="2" x2="22" y2="22" /><path
										d="M17.547 17.547a8.553 8.553 0 0 1-5.547 1.953 8.8 8.8 0 0 1-8.547-5.5 10.87 10.87 0 0 1 1.761-3.239"
									/><path
										d="M9.88 4.22a8.8 8.8 0 0 1 1.62-.22 8.82 8.82 0 0 1 8.547 5.5 10.64 10.64 0 0 1-1.341 2.871"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{:else}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<path
										d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0z"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Confirm Password -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-reset-confirm-password"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						CONFIRM NEW PASSWORD
					</label>
					<div class="relative flex items-center">
						<span class="absolute left-3.5 pointer-events-none text-slate-400">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								style="margin-left: 10px"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2.2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
								/>
							</svg>
						</span>
						<input
							id="{formId}-reset-confirm-password"
							type={showForgotConfirmPassword ? 'text' : 'password'}
							bind:value={forgotConfirmPassword}
							placeholder="Confirm new password"
							class="w-full pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							style="padding-left: 40px;"
							required
						/>
						<button
							type="button"
							onclick={() => (showForgotConfirmPassword = !showForgotConfirmPassword)}
							class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none"
							aria-label={showForgotConfirmPassword ? 'Hide password' : 'Show password'}
						>
							{#if showForgotConfirmPassword}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<line x1="2" y1="2" x2="22" y2="22" /><path
										d="M17.547 17.547a8.553 8.553 0 0 1-5.547 1.953 8.8 8.8 0 0 1-8.547-5.5 10.87 10.87 0 0 1 1.761-3.239"
									/><path
										d="M9.88 4.22a8.8 8.8 0 0 1 1.62-.22 8.82 8.82 0 0 1 8.547 5.5 10.64 10.64 0 0 1-1.341 2.871"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{:else}
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<path
										d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0z"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Password Requirements Checklist -->
				{#if forgotNewPassword}
					<div
						class="p-4 bg-slate-50 border border-border-base rounded-xl flex flex-col gap-2.5"
						transition:slide={{ duration: 200 }}
					>
						<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
							Password Requirements
						</div>
						<div class="flex flex-col gap-2">
							<div class="flex items-center gap-2 text-xs">
								<span class={resetHasMinLength ? 'text-emerald-600 font-bold' : 'text-slate-400'}>
									{resetHasMinLength ? '✓' : '○'}
								</span>
								<span class={resetHasMinLength ? 'text-slate-700' : 'text-slate-500'}
									>Minimum 8 characters</span
								>
							</div>
							<div class="flex items-center gap-2 text-xs">
								<span class={resetHasUppercase ? 'text-emerald-600 font-bold' : 'text-slate-400'}>
									{resetHasUppercase ? '✓' : '○'}
								</span>
								<span class={resetHasUppercase ? 'text-slate-700' : 'text-slate-500'}
									>One uppercase letter</span
								>
							</div>
							<div class="flex items-center gap-2 text-xs">
								<span class={resetHasLowercase ? 'text-emerald-600 font-bold' : 'text-slate-400'}>
									{resetHasLowercase ? '✓' : '○'}
								</span>
								<span class={resetHasLowercase ? 'text-slate-700' : 'text-slate-500'}
									>One lowercase letter</span
								>
							</div>
							<div class="flex items-center gap-2 text-xs">
								<span class={resetHasNumber ? 'text-emerald-600 font-bold' : 'text-slate-400'}>
									{resetHasNumber ? '✓' : '○'}
								</span>
								<span class={resetHasNumber ? 'text-slate-700' : 'text-slate-500'}>One number</span>
							</div>
							<div class="flex items-center gap-2 text-xs">
								<span class={resetHasSpecialChar ? 'text-emerald-600 font-bold' : 'text-slate-400'}>
									{resetHasSpecialChar ? '✓' : '○'}
								</span>
								<span class={resetHasSpecialChar ? 'text-slate-700' : 'text-slate-500'}
									>One special character</span
								>
							</div>
						</div>
					</div>
				{/if}

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={submitting ||
						!resetAllPasswordRequirementsMet ||
						!resetPasswordsMatch ||
						forgotOtp.trim() === ''}
					class="w-full py-3.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-widest uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
				>
					{#if submitting}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						RESETTING...
					{:else}
						RESET PASSWORD
					{/if}
				</button>

				<div class="flex justify-center items-center text-xs mt-2">
					<button
						type="button"
						onclick={() => {
							viewState = 'forgot';
							errorMsg = '';
						}}
						class="text-slate-500 hover:text-slate-800 font-semibold transition duration-200"
					>
						← Back to Email Input
					</button>
				</div>
			</form>
		{:else}
			<!-- Login Form View -->
			<div class="p-6 sm:p-8 border-b border-border-base bg-slate-50/50">
				<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
					STUDENT PORTAL
				</div>
				<h2 class="text-2xl font-bold text-inst-navy font-serif leading-tight mt-1">
					Welcome Back!
				</h2>
				<p class="text-slate-500 text-xs mt-1">
					Sign in to continue to your iSPARC student account.
				</p>
			</div>

			<form onsubmit={handleSubmit} class="p-6 sm:p-8 flex flex-col gap-6">
				{#if errorMsg}
					<div
						class="p-3.5 bg-rose-50 border border-rose-200 text-rose-700 text-xs font-semibold rounded-lg"
						transition:slide={{ duration: 200 }}
					>
						{errorMsg}
					</div>
				{/if}

				<!-- Institutional Email -->
				<div class="flex flex-col gap-1.5">
					<label for="{formId}-email" class="text-[11px] font-bold text-slate-700 tracking-wider">
						EMAIL
					</label>
					<div class="relative flex items-center">
						<span class="absolute left-3.5 pointer-events-none text-slate-400">
							<!-- Envelope SVG -->
							<svg
								xmlns="http://www.w3.org/2000/svg"
								style="margin-left: 10px"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2.2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75"
								/>
							</svg>
						</span>
						<input
							id="{formId}-email"
							type="email"
							bind:value={email}
							placeholder="student.iips@gmail.com"
							class="w-full pr-3.5 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							style="padding-left: 40px;"
							required
						/>
					</div>
				</div>

				<!-- Captcha Section -->
				<div class="flex flex-col gap-1.5">
					<label for="{formId}-captcha" class="text-[11px] font-bold text-slate-700 tracking-wider">
						CAPTCHA SECURITY <span class="text-acad-red">*</span>
					</label>
					<div class="flex gap-2 items-center">
						<!-- Captcha display -->
						<div
							class="flex-grow flex items-center justify-center bg-slate-100 border border-slate-200 rounded-lg py-2.5 px-3 font-mono font-bold text-slate-800 tracking-wider select-none text-[15px] min-h-[42px]"
						>
							{#if captchaLoading}
								<svg class="animate-spin h-5 w-5 text-slate-500" fill="none" viewBox="0 0 24 24">
									<circle
										class="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										stroke-width="4"
									></circle>
									<path
										class="opacity-75"
										fill="currentColor"
										d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
									></path>
								</svg>
							{:else}
								{captchaQuestion || 'No captcha loaded'}
							{/if}
						</div>

						<!-- Refresh button -->
						<button
							type="button"
							onclick={fetchCaptcha}
							disabled={captchaLoading}
							class="shrink-0 p-2.5 bg-slate-50 hover:bg-slate-100 border border-slate-200 rounded-lg text-slate-500 hover:text-slate-800 transition duration-200"
							title="Refresh Captcha"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2.2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="w-4 h-4"
							>
								<path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
								<path d="M3 3v5h5" />
								<path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16" />
								<path d="M16 16h5v5" />
							</svg>
						</button>

						<!-- Captcha Input -->
						<input
							id="{formId}-captcha"
							type="text"
							bind:value={captchaAnswer}
							placeholder="Answer"
							class="w-24 px-3 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-center text-slate-800 font-bold focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							required
						/>
					</div>
				</div>

				<!-- Password -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-password"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						PASSWORD
					</label>
					<div class="relative flex items-center">
						<span class="absolute left-3.5 pointer-events-none text-slate-400">
							<!-- Lock SVG -->
							<svg
								xmlns="http://www.w3.org/2000/svg"
								style="margin-left: 10px"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2.2"
								stroke="currentColor"
								class="w-4 h-4"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
								/>
							</svg>
						</span>
						<input
							id="{formId}-password"
							type={showPassword ? 'text' : 'password'}
							bind:value={password}
							placeholder="Enter your password"
							class="w-full pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200"
							style="padding-left: 40px;"
							required
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 focus:outline-none"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showPassword}
								<!-- Eye-off SVG -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<line x1="2" y1="2" x2="22" y2="22" /><path
										d="M17.547 17.547a8.553 8.553 0 0 1-5.547 1.953 8.8 8.8 0 0 1-8.547-5.5 10.87 10.87 0 0 1 1.761-3.239"
									/><path
										d="M9.88 4.22a8.8 8.8 0 0 1 1.62-.22 8.82 8.82 0 0 1 8.547 5.5 10.64 10.64 0 0 1-1.341 2.871"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{:else}
								<!-- Eye SVG -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<path
										d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0z"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Remember Me & Forgot Password -->
				<div class="flex items-center justify-between text-xs">
					<label class="flex items-center gap-2 cursor-pointer select-none">
						<input
							type="checkbox"
							bind:checked={rememberMe}
							class="w-3.5 h-3.5 rounded border-border-base text-inst-navy accent-inst-navy cursor-pointer"
						/>
						<span class="text-slate-500 font-medium hover:text-slate-700 transition-colors"
							>Remember Me</span
						>
					</label>

					<button
						type="button"
						onclick={() => {
							viewState = 'forgot';
							errorMsg = '';
						}}
						class="font-bold text-acad-red hover:underline transition-colors focus:outline-none"
					>
						Forgot Password?
					</button>
				</div>

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={submitting || !isFormValid}
					class="w-full py-3.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-widest uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
				>
					{#if submitting}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						LOGGING IN...
					{:else}
						LOGIN
					{/if}
				</button>

				<!-- Divider -->
				<div class="relative flex py-2 items-center">
					<div class="flex-grow border-t border-border-base"></div>
					<span
						class="flex-shrink mx-4 text-[10px] font-bold text-slate-400 uppercase tracking-widest"
						>OR</span
					>
					<div class="flex-grow border-t border-border-base"></div>
				</div>

				<!-- Register Redirect Link -->
				<div class="flex flex-col gap-3 text-center">
					<div class="text-xs text-slate-500">Don't have an account?</div>
					<a
						href="/register"
						class="w-full py-3 bg-white border border-inst-navy hover:bg-inst-navy hover:text-white text-inst-navy font-semibold text-xs tracking-wider uppercase rounded-xl transition duration-200 text-center"
					>
						Create Student Account
					</a>
				</div>

				<!-- Help Section notice -->
				<div class="mt-2.5 p-4 bg-slate-50 border border-border-base rounded-xl">
					<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-2">
						NEED HELP?
					</div>
					<ul class="flex flex-col gap-2 list-none">
						{#each helpItems as item}
							<li class="flex items-center gap-2 text-xs text-slate-500">
								<span class="w-1.5 h-1.5 bg-slate-400 rounded-full shrink-0"></span>
								<span>{item}</span>
							</li>
						{/each}
					</ul>
					<div class="mt-3">
						<button
							type="button"
							class="w-full py-2 bg-white border border-border-base text-[11px] font-semibold text-slate-650 hover:bg-slate-50 rounded-lg transition duration-200"
						>
							Contact iSPARC Support
						</button>
					</div>
				</div>
			</form>
		{/if}
	</main>
</section>

<style>
	.animate-fade-in {
		opacity: 0;
		transform: translateY(-10px);
		animation: fadeIn 0.8s ease-out forwards;
	}

	.animate-fade-up {
		opacity: 0;
		transform: translateY(25px);
		animation: fadeUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
		animation-delay: 0.15s;
	}

	@keyframes fadeIn {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes fadeUp {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
