<script lang="ts">
	// Define props
	let { onUploadCertificate }: { onUploadCertificate?: () => void } = $props();

	// Types
	interface RegisteredActivity {
		name: string;
		category: string;
		enrollmentDate: string;
		activityDate: string;
		credits: number;
		status: 'Registered' | 'Completed' | 'Pending Verification' | 'Verified' | 'Rejected';
		action: 'View Details' | 'View Certificate' | 'Upload Certificate';
	}

	interface JourneyStep {
		label: string;
		subtext: string;
		status: 'completed' | 'current' | 'pending';
	}

	interface UpcomingActivity {
		name: string;
		date: string;
		time: string;
		venue: string;
		credits: number;
		daysLeft: string;
	}

	interface CreditCategorySummary {
		category: string;
		creditsEarned: number;
		creditsRequired: number;
		percentage: number;
	}

	interface RecentUpdate {
		type: 'success' | 'star' | 'doc' | 'db';
		text: string;
		time: string;
	}

	// 7 mock records exactly matching the screenshot
	const registeredActivities: RegisteredActivity[] = [
		{
			name: 'National Hackathon 2026',
			category: 'Technical',
			enrollmentDate: '10 Jun 2026',
			activityDate: '05 Jul 2026',
			credits: 15,
			status: 'Registered',
			action: 'View Details'
		},
		{
			name: 'Inter College Debate',
			category: 'Public Speaking',
			enrollmentDate: '05 Jun 2026',
			activityDate: '20 Jun 2026',
			credits: 10,
			status: 'Completed',
			action: 'View Certificate'
		},
		{
			name: 'Blood Donation Drive',
			category: 'Social Service',
			enrollmentDate: '15 May 2026',
			activityDate: '25 May 2026',
			credits: 5,
			status: 'Pending Verification',
			action: 'Upload Certificate'
		},
		{
			name: 'Robotics Workshop',
			category: 'Technical',
			enrollmentDate: '02 Apr 2026',
			activityDate: '18 Apr 2026',
			credits: 12,
			status: 'Verified',
			action: 'View Certificate'
		},
		{
			name: 'Annual Cultural Fest',
			category: 'Cultural',
			enrollmentDate: '01 Apr 2026',
			activityDate: '20 Mar 2026',
			credits: 8,
			status: 'Completed',
			action: 'View Certificate'
		},
		{
			name: 'Research Presentation Contest',
			category: 'Academic',
			enrollmentDate: '15 Feb 2026',
			activityDate: '10 Mar 2026',
			credits: 12,
			status: 'Verified',
			action: 'View Certificate'
		},
		{
			name: 'National Science Olympiad',
			category: 'Technical',
			enrollmentDate: '30 Jan 2026',
			activityDate: '05 Feb 2026',
			credits: 10,
			status: 'Rejected',
			action: 'View Details'
		}
	];

	// Upcoming activities
	const upcomingActivities: UpcomingActivity[] = [
		{
			name: 'National Hackathon 2026',
			date: '5 July 2026',
			time: '10:00 AM',
			venue: 'IIPS Auditorium',
			credits: 15,
			daysLeft: '12d left'
		},
		{
			name: 'Innovation Challenge',
			date: '12 July 2026',
			time: '09:00 AM',
			venue: 'Seminar Hall B',
			credits: 10,
			daysLeft: '19d left'
		},
		{
			name: 'Technical Workshop Series',
			date: '20 July 2026',
			time: '11:00 AM',
			venue: 'Computer Lab 3',
			credits: 8,
			daysLeft: '27d left'
		},
		{
			name: 'Research Symposium',
			date: '2 Aug 2026',
			time: '10:30 AM',
			venue: 'Conference Room',
			credits: 12,
			daysLeft: '40d left'
		}
	];

	// Credit Summary category values
	const creditCategorySummaries: CreditCategorySummary[] = [
		{
			category: 'Technical Activities',
			creditsEarned: 48,
			creditsRequired: 60,
			percentage: 80
		},
		{
			category: 'Public Speaking',
			creditsEarned: 24,
			creditsRequired: 40,
			percentage: 60
		},
		{
			category: 'Social Service',
			creditsEarned: 18,
			creditsRequired: 30,
			percentage: 60
		},
		{
			category: 'Sports',
			creditsEarned: 12,
			creditsRequired: 30,
			percentage: 40
		}
	];

	// Recent updates timeline data
	const recentUpdates: RecentUpdate[] = [
		{
			type: 'success',
			text: 'National Hackathon registration confirmed',
			time: '2 days ago'
		},
		{
			type: 'star',
			text: '5 credits added for Robotics Workshop',
			time: '5 days ago'
		},
		{
			type: 'doc',
			text: 'Certificate verified for Annual Fest',
			time: '1 week ago'
		},
		{
			type: 'db',
			text: 'Enrolled in Innovation Challenge',
			time: '2 weeks ago'
		}
	];

	// Selection state: selected row index (default to index 1 "Inter College Debate" as shown in the screenshot)
	let selectedActivityIndex = $state(1);

	// Get selected activity details
	let selectedActivity = $derived(registeredActivities[selectedActivityIndex]);

	// Generate tracking journey steps dynamically based on selected activity status
	let journeySteps = $derived.by<JourneyStep[]>(() => {
		const status = selectedActivity.status;
		const baseSteps = [
			{ label: 'Enrolled', subtext: 'Registered for the activity' },
			{ label: 'Activity Started', subtext: 'Event began on scheduled date' },
			{ label: 'Completed', subtext: 'Participated in all sessions' },
			{ label: 'Certificate Submitted', subtext: 'Uploaded proof of participation' },
			{ label: 'Verified', subtext: 'Awaiting mentor verification' },
			{ label: 'Credits Awarded', subtext: 'Credits will reflect in profile' }
		];

		if (status === 'Registered') {
			return [
				{ ...baseSteps[0], status: 'completed' },
				{ ...baseSteps[1], status: 'current' },
				{ ...baseSteps[2], status: 'pending' },
				{ ...baseSteps[3], status: 'pending' },
				{ ...baseSteps[4], status: 'pending' },
				{ ...baseSteps[5], status: 'pending' }
			];
		} else if (status === 'Completed') {
			return [
				{ ...baseSteps[0], status: 'completed' },
				{ ...baseSteps[1], status: 'completed' },
				{ ...baseSteps[2], status: 'completed' },
				{ ...baseSteps[3], status: 'current' },
				{ ...baseSteps[4], status: 'pending' },
				{ ...baseSteps[5], status: 'pending' }
			];
		} else if (status === 'Pending Verification') {
			return [
				{ ...baseSteps[0], status: 'completed' },
				{ ...baseSteps[1], status: 'completed' },
				{ ...baseSteps[2], status: 'completed' },
				{ ...baseSteps[3], status: 'completed' },
				{ ...baseSteps[4], status: 'current' },
				{ ...baseSteps[5], status: 'pending' }
			];
		} else if (status === 'Verified') {
			return [
				{ ...baseSteps[0], status: 'completed' },
				{ ...baseSteps[1], status: 'completed' },
				{ ...baseSteps[2], status: 'completed' },
				{ ...baseSteps[3], status: 'completed' },
				{ ...baseSteps[4], status: 'completed' },
				{ ...baseSteps[5], status: 'current' }
			];
		} else if (status === 'Rejected') {
			return [
				{ ...baseSteps[0], status: 'completed' },
				{ ...baseSteps[1], status: 'completed' },
				{ ...baseSteps[2], status: 'completed' },
				{ ...baseSteps[3], status: 'completed' },
				{
					label: 'Verification Rejected',
					subtext: 'Certificate rejected during audit',
					status: 'current' as const
				},
				{ ...baseSteps[5], status: 'pending' }
			];
		}

		return baseSteps.map((step) => ({ ...step, status: 'pending' as const }));
	});

	function handleActionClick(act: RegisteredActivity) {
		if (act.action === 'Upload Certificate') {
			if (onUploadCertificate) {
				onUploadCertificate();
			} else {
				alert('Opening Upload Certificate modal...');
			}
		} else {
			alert(`Action: "${act.action}" triggered for "${act.name}"`);
		}
	}
</script>

<div class="space-y-6">
	<!-- ==================== 1. KPI WIDGETS ROW ==================== -->
	<section
		class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4"
		aria-label="KPI Summary Cards"
	>
		<!-- KPI Card 1: Total Enrollments -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex items-center gap-4.5 shadow-xs"
		>
			<div
				class="w-12 h-12 rounded-lg bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-500 shrink-0"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
					/>
				</svg>
			</div>
			<div class="flex flex-col">
				<span class="text-2xl font-bold font-serif text-slate-900 leading-tight">12</span>
				<span class="text-xs font-bold text-slate-400 mt-0.5 tracking-wide uppercase"
					>Total Enrollments</span
				>
			</div>
		</div>

		<!-- KPI Card 2: Upcoming Activities -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex items-center gap-4.5 shadow-xs"
		>
			<div
				class="w-12 h-12 rounded-lg bg-blue-50 border border-blue-100 flex items-center justify-center text-blue-600 shrink-0"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5m-9-6h.008v.008H12v-.008zM12 15h.008v.008H12V15zm0 2.25h.008v.008H12v-.008zM9.75 15h.008v.008H9.75V15zm0 2.25h.008v.008H9.75v-.008zM7.5 15h.008v.008H7.5V15zm0 2.25h.008v.008H7.5v-.008zm6.75-4.5h.008v.008h-.008v-.008zm0 2.25h.008v.008h-.008V15zm0 2.25h.008v.008h-.008v-.008zm2.25-4.5h.008v.008H16.5v-.008zm0 2.25h.008v.008H16.5V15z"
					/>
				</svg>
			</div>
			<div class="flex flex-col">
				<span class="text-2xl font-bold font-serif text-slate-900 leading-tight">4</span>
				<span class="text-xs font-bold text-slate-400 mt-0.5 tracking-wide uppercase"
					>Upcoming Activities</span
				>
			</div>
		</div>

		<!-- KPI Card 3: Completed Activities -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex items-center gap-4.5 shadow-xs"
		>
			<div
				class="w-12 h-12 rounded-lg bg-emerald-50 border border-emerald-100 flex items-center justify-center text-emerald-600 shrink-0"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.599-3.751A11.959 11.959 0 0112 2.714z"
					/>
				</svg>
			</div>
			<div class="flex flex-col">
				<span class="text-2xl font-bold font-serif text-slate-900 leading-tight">6</span>
				<span class="text-xs font-bold text-slate-400 mt-0.5 tracking-wide uppercase"
					>Completed Activities</span
				>
			</div>
		</div>

		<!-- KPI Card 4: Pending Verification -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex items-center gap-4.5 shadow-xs"
		>
			<div
				class="w-12 h-12 rounded-lg bg-amber-50 border border-amber-100 flex items-center justify-center text-amber-600 shrink-0"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-6 h-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			</div>
			<div class="flex flex-col">
				<span class="text-2xl font-bold font-serif text-slate-900 leading-tight">2</span>
				<span class="text-xs font-bold text-slate-400 mt-0.5 tracking-wide uppercase"
					>Pending Verification</span
				>
			</div>
		</div>
	</section>

	<!-- ==================== 2. REGISTERED ACTIVITIES TABLE ==================== -->
	<section class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs flex flex-col">
		<div class="flex items-center justify-between pb-4 border-b border-slate-100">
			<h2 class="text-base font-bold font-serif text-[#0B1535]">My Registered Activities</h2>
			<span class="text-xs text-slate-405 font-bold tracking-wide uppercase"
				>{registeredActivities.length} records</span
			>
		</div>

		<!-- Desktop Table -->
		<div class="overflow-x-auto mt-4">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr
						class="text-[10px] font-bold text-[#6B7280] uppercase tracking-widest border-b border-slate-100 bg-slate-50/50"
					>
						<th class="py-3 px-4">Activity Name</th>
						<th class="py-3 px-4">Category</th>
						<th class="py-3 px-4">Enrollment Date</th>
						<th class="py-3 px-4">Activity Date</th>
						<th class="py-3 px-4 text-center">Credits</th>
						<th class="py-3 px-4 text-center">Status</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100">
					{#each registeredActivities as act, idx}
						<tr
							onclick={() => (selectedActivityIndex = idx)}
							class="hover:bg-slate-50/80 cursor-pointer transition-colors duration-150 {selectedActivityIndex ===
							idx
								? 'bg-[#881B1B]/5 font-semibold'
								: ''}"
						>
							<td
								class="py-3.5 px-4 font-bold text-[#0B1535] leading-relaxed max-w-[220px] truncate"
							>
								{act.name}
							</td>
							<td class="py-3.5 px-4 text-slate-600 font-medium">
								{act.category}
							</td>
							<td class="py-3.5 px-4 text-slate-500 whitespace-nowrap">
								{act.enrollmentDate}
							</td>
							<td class="py-3.5 px-4 text-slate-500 whitespace-nowrap">
								{act.activityDate}
							</td>
							<td class="py-3.5 px-4 text-center font-extrabold text-[#881B1B]">
								{act.credits}
							</td>
							<td class="py-3.5 px-4 text-center whitespace-nowrap">
								<span
									class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-sm text-[10px] font-extrabold tracking-wide uppercase border
									{act.status === 'Registered' ? 'bg-blue-50 text-blue-700 border-blue-200' : ''}
									{act.status === 'Completed' || act.status === 'Verified'
										? 'bg-emerald-50 text-emerald-700 border-emerald-200'
										: ''}
									{act.status === 'Pending Verification' ? 'bg-amber-50 text-amber-700 border-amber-200' : ''}
									{act.status === 'Rejected' ? 'bg-rose-50 text-rose-700 border-rose-200' : ''}"
								>
									<span
										class="w-1.5 h-1.5 rounded-full
										{act.status === 'Registered' ? 'bg-blue-500' : ''}
										{act.status === 'Completed' || act.status === 'Verified' ? 'bg-emerald-500' : ''}
										{act.status === 'Pending Verification' ? 'bg-amber-500' : ''}
										{act.status === 'Rejected' ? 'bg-rose-500' : ''}"
									></span>
									{act.status}
								</span>
							</td>
							<td
								class="py-3.5 px-4 text-right whitespace-nowrap"
								onclick={(e) => e.stopPropagation()}
							>
								<button
									onclick={() => handleActionClick(act)}
									class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[11px] font-bold border transition-colors duration-150 cursor-pointer
										{act.action === 'Upload Certificate'
										? 'border-amber-200 bg-amber-50 text-amber-700 hover:bg-amber-100'
										: 'border-slate-200 bg-white hover:bg-slate-50 text-slate-700'}"
								>
									{#if act.action === 'Upload Certificate'}
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.5"
											stroke="currentColor"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"
											/>
										</svg>
									{:else if act.action === 'View Certificate'}
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.2"
											stroke="currentColor"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
											/>
										</svg>
									{:else}
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.2"
											stroke="currentColor"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
											/>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
											/>
										</svg>
									{/if}
									{act.action}
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</section>

	<!-- ==================== 3. DOUBLE COLUMN LAYOUT ==================== -->
	<section class="grid grid-cols-1 lg:grid-cols-12 gap-6">
		<!-- LEFT COLUMN (Journey Tracker & Upcoming list) -->
		<div class="lg:col-span-8 space-y-6 flex flex-col">
			<!-- A. Participation Journey -->
			<div
				class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs flex-grow flex flex-col"
			>
				<div class="flex items-start justify-between pb-4 border-b border-slate-100 mb-5">
					<div>
						<h2 class="text-base font-bold font-serif text-[#0B1535]">Participation Journey</h2>
						<p class="text-[10.5px] font-bold text-slate-405 mt-1 uppercase tracking-wide">
							{selectedActivity.name} - Tracking Status
						</p>
					</div>
					<span
						class="px-2 py-0.5 bg-[#881B1B]/10 text-[#881B1B] text-[10px] font-extrabold uppercase rounded-sm tracking-wide"
					>
						Active Tracking
					</span>
				</div>

				<!-- Vertical Timeline -->
				<div class="space-y-6 relative pl-3.5 my-2 flex-grow flex flex-col justify-center">
					<!-- Vertical Line connector -->
					<div class="absolute left-6.5 top-2.5 bottom-2.5 w-0.5 bg-slate-100"></div>

					{#each journeySteps as step}
						<div class="flex items-start gap-4 relative z-10">
							<!-- Timeline Indicator Dot -->
							<div
								class="flex items-center justify-center shrink-0 w-6 h-6 rounded-full border bg-white mt-0.5
								{step.status === 'completed' ? 'border-[#0B1535] bg-[#0B1535] text-white' : ''}
								{step.status === 'current' ? 'border-2 border-[#0B1535] bg-white text-[#0B1535]' : ''}
								{step.status === 'pending' ? 'border-slate-200 bg-white text-slate-300' : ''}"
							>
								{#if step.status === 'completed'}
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="3"
										stroke="currentColor"
										class="w-3.5 h-3.5"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M4.5 12.75l6 6 9-13.5"
										/>
									</svg>
								{:else if step.status === 'current'}
									<span class="w-1.5 h-1.5 bg-[#0B1535] rounded-full animate-ping"></span>
								{:else}
									<span class="w-1 h-1 bg-slate-200 rounded-full"></span>
								{/if}
							</div>

							<!-- Timeline Label -->
							<div class="flex-grow flex items-start justify-between min-w-0">
								<div>
									<h4
										class="text-xs font-bold leading-tight
										{step.status === 'completed' || step.status === 'current' ? 'text-[#0B1535]' : 'text-slate-400'}"
									>
										{step.label}
									</h4>
									<p class="text-[10px] font-medium text-slate-400 mt-0.5 leading-relaxed truncate">
										{step.subtext}
									</p>
								</div>

								{#if step.status === 'current'}
									<span
										class="shrink-0 px-2 py-0.5 bg-amber-50 text-amber-700 border border-amber-200/50 text-[8px] font-extrabold uppercase rounded-xs tracking-wider"
									>
										Current
									</span>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- B. Upcoming Activities List -->
			<div class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs">
				<div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-4">
					<h2 class="text-base font-bold font-serif text-[#0B1535]">Upcoming Activities</h2>
					<span
						class="px-2.5 py-1 bg-blue-50 text-blue-700 border border-blue-100 rounded-lg text-[9px] font-extrabold uppercase tracking-wide"
					>
						Active Registration
					</span>
				</div>

				<div class="divide-y divide-slate-100">
					{#each upcomingActivities as activity}
						<div class="py-3.5 flex items-center justify-between gap-4 first:pt-0 last:pb-0">
							<div class="flex items-center gap-3.5 min-w-0">
								<!-- Icon -->
								<div
									class="w-10 h-10 rounded-lg bg-blue-50/80 border border-blue-100/50 flex items-center justify-center text-blue-600 shrink-0"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2"
										stroke="currentColor"
										class="w-5 h-5"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5"
										/>
									</svg>
								</div>

								<div class="min-w-0">
									<h4 class="text-xs font-bold text-slate-800 leading-tight truncate">
										{activity.name}
									</h4>
									<div
										class="flex flex-wrap items-center gap-x-2 gap-y-0.5 mt-1 text-[10px] text-slate-405 font-semibold font-sans"
									>
										<span class="flex items-center gap-1">
											<svg
												xmlns="http://www.w3.org/2000/svg"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="2"
												stroke="currentColor"
												class="w-3 h-3 text-slate-400"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
												/>
											</svg>
											{activity.date} @ {activity.time}
										</span>
										<span class="text-slate-300">•</span>
										<span class="flex items-center gap-1">
											<svg
												xmlns="http://www.w3.org/2000/svg"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="2"
												stroke="currentColor"
												class="w-3 h-3 text-slate-400"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
												/>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
												/>
											</svg>
											{activity.venue}
										</span>
									</div>
								</div>
							</div>

							<div class="text-right shrink-0 flex flex-col items-end">
								<span class="text-xs font-extrabold text-[#881B1B]">{activity.credits} Credits</span
								>
								<span
									class="text-[9.5px] font-bold text-slate-400 block mt-1 uppercase tracking-wide"
									>{activity.daysLeft}</span
								>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>

		<!-- RIGHT COLUMN (Credit Summary & Updates list) -->
		<div class="lg:col-span-4 space-y-6 flex flex-col">
			<!-- A. Enrollment Credit Summary -->
			<div
				class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs flex flex-col justify-between"
			>
				<div>
					<h2 class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100">
						Enrollment Credit Summary
					</h2>

					<!-- Category progress list -->
					<div class="space-y-4 my-5">
						{#each creditCategorySummaries as cat}
							<div class="space-y-1.5">
								<div class="flex justify-between text-xs font-bold text-slate-800">
									<span class="text-[10px] font-extrabold text-slate-405 uppercase tracking-wide"
										>{cat.category}</span
									>
									<span class="text-slate-900">{cat.creditsEarned} Credits</span>
								</div>
								<div class="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
									<div
										class="h-full bg-[#0B1535] rounded-full"
										style="width: {cat.percentage}%"
									></div>
								</div>
							</div>
						{/each}
					</div>
				</div>

				<div class="border-t border-slate-150 pt-5 mt-2">
					<div class="flex justify-between text-xs font-extrabold text-slate-800 mb-2">
						<span class="text-[11px] font-extrabold text-[#0B1535] uppercase tracking-wide"
							>Total Credits vs Required</span
						>
						<span class="text-[#881B1B]">118 / 200</span>
					</div>
					<div class="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
						<div class="h-full bg-[#881B1B] rounded-full" style="width: 59%"></div>
					</div>
					<div
						class="flex justify-end text-[10px] font-extrabold text-[#881B1B] mt-1.5 uppercase tracking-wide"
					>
						59% Complete
					</div>
				</div>
			</div>

			<!-- B. Recent Enrollment Updates -->
			<div
				class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs flex-grow flex flex-col"
			>
				<h2
					class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100 mb-4"
				>
					Recent Enrollment Updates
				</h2>

				<div class="space-y-4.5 my-2 flex-grow flex flex-col justify-center">
					{#each recentUpdates as update}
						<div class="flex items-start gap-3">
							<!-- Circle status indicator icon -->
							<div
								class="w-8 h-8 rounded-full border flex items-center justify-center shrink-0
								{update.type === 'success' ? 'bg-emerald-50 border-emerald-100 text-emerald-600' : ''}
								{update.type === 'star' ? 'bg-blue-50 border-blue-100 text-blue-600' : ''}
								{update.type === 'doc' ? 'bg-emerald-50 border-emerald-100 text-emerald-600' : ''}
								{update.type === 'db' ? 'bg-blue-50 border-blue-100 text-blue-600' : ''}"
							>
								{#if update.type === 'success'}
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2.5"
										stroke="currentColor"
										class="w-4 h-4"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
										/>
									</svg>
								{:else if update.type === 'star'}
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2.2"
										stroke="currentColor"
										class="w-4 h-4"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M11.48 3.499c.174-.627 1.06-.627 1.234 0l1.63 4.887a1 1 0 00.95.69h5.132c.656 0 .928.843.429 1.26l-4.153 3.493a1 1 0 00-.342 1.053l1.63 4.887c.174.627-.587 1.181-1.129.782l-4.154-3.493a1 1 0 00-1.17 0l-4.154 3.493c-.542.399-1.303-.155-1.129-.782l1.63-4.887a1 1 0 00-.342-1.053L2.42 10.336c-.498-.417-.226-1.26.429-1.26H7.98a1 1 0 00.95-.69l1.63-4.887z"
										/>
									</svg>
								{:else if update.type === 'doc'}
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2.2"
										stroke="currentColor"
										class="w-4 h-4"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
										/>
									</svg>
								{:else}
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2.2"
										stroke="currentColor"
										class="w-4 h-4"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125v-3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75"
										/>
									</svg>
								{/if}
							</div>

							<div class="min-w-0">
								<p class="text-xs text-slate-700 font-semibold leading-relaxed">{update.text}</p>
								<span class="text-[10px] text-slate-400 block mt-0.5">{update.time}</span>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</section>
</div>
