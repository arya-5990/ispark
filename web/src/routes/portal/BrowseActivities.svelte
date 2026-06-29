<script lang="ts">
	import { slide } from 'svelte/transition';

	// Define interfaces for activities
	interface Activity {
		id: number;
		name: string;
		category: string;
		description: string;
		credits: number;
		mode: string;
		regDeadline: string;
		activityDate: string;
		venue: string;
		coordinator: string;
		status: 'Open' | 'Closed' | 'Closing Soon';
		daysLeft?: string;
	}

	// Mock list of activities
	const allActivities: Activity[] = [
		{
			id: 1,
			name: 'National Hackathon 2026',
			category: 'TECHNICAL',
			description:
				'A 36-hour coding challenge open to all undergraduate students. Build innovative solutions for real-world problems across domains.',
			credits: 15,
			mode: 'Offline',
			regDeadline: '30 June 2026',
			activityDate: '5 July 2026',
			venue: 'IIPS Auditorium',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Closing Soon',
			daysLeft: '3d left'
		},
		{
			id: 2,
			name: 'Inter-College Athletics Meet',
			category: 'SPORTS',
			description:
				'Annual inter-college athletics championship. Compete in track and field events representing IIPS at the university level.',
			credits: 10,
			mode: 'Offline',
			regDeadline: '25 June 2026',
			activityDate: '8 July 2026',
			venue: 'DAVV Sports Ground',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Open',
			daysLeft: '7d left'
		},
		{
			id: 3,
			name: 'National Science Olympiad',
			category: 'RESEARCH',
			description:
				'Prestigious national-level science competition covering physics, chemistry, and biology. Written exam followed by practical rounds.',
			credits: 20,
			mode: 'Hybrid',
			regDeadline: '15 July 2026',
			activityDate: '20 July 2026',
			venue: 'IIPS Seminar Hall',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Open'
		},
		{
			id: 4,
			name: 'Inter College Debate Championship',
			category: 'PUBLIC SPEAKING',
			description:
				'Parliamentary-style debate on contemporary socio-political topics. Individual and team participation categories available.',
			credits: 12,
			mode: 'Offline',
			regDeadline: '28 June 2026',
			activityDate: '3 July 2026',
			venue: 'IIPS Conference Hall',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Closing Soon',
			daysLeft: '5d left'
		},
		{
			id: 5,
			name: 'Blood Donation Camp',
			category: 'SOCIAL SERVICE',
			description:
				'Community health initiative in partnership with District Hospital Indore. Volunteers earn certified social service credit.',
			credits: 8,
			mode: 'Offline',
			regDeadline: '22 June 2026',
			activityDate: '1 July 2026',
			venue: 'IIPS Main Ground',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Closed',
			daysLeft: '9d left'
		},
		{
			id: 6,
			name: 'Annual Cultural Fest — Dance',
			category: 'CULTURAL',
			description:
				'Classical and contemporary dance competition as part of the Annual Cultural Festival. Solo and group categories.',
			credits: 10,
			mode: 'Offline',
			regDeadline: '10 July 2026',
			activityDate: '18 July 2026',
			venue: 'Open Air Theatre, DAVV',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Open'
		},
		{
			id: 7,
			name: 'AI & Machine Learning Bootcamp',
			category: 'TECHNICAL',
			description:
				'Hands-on bootcamp on deep learning, generative AI models, and model tuning. Earn certificates and practical project credits.',
			credits: 15,
			mode: 'Online',
			regDeadline: '12 July 2026',
			activityDate: '16 July 2026',
			venue: 'Google Meet',
			coordinator: 'Dr. Sanjay Tanwani',
			status: 'Open'
		},
		{
			id: 8,
			name: 'Swachh Bharat Cleanliness Drive',
			category: 'SOCIAL SERVICE',
			description:
				'Campus-wide cleanliness and awareness drive. Volunteer to help make IIPS and DAVV plastic-free zones.',
			credits: 6,
			mode: 'Offline',
			regDeadline: '05 July 2026',
			activityDate: '07 July 2026',
			venue: 'DAVV Campus',
			coordinator: 'Dr. K. K. Pandey',
			status: 'Open'
		},
		{
			id: 9,
			name: 'NSS Youth Leadership Summit',
			category: 'LEADERSHIP',
			description:
				'A leadership workshop teaching communication, organizing skills, and social responsibility principles.',
			credits: 10,
			mode: 'Offline',
			regDeadline: '14 July 2026',
			activityDate: '18 July 2026',
			venue: 'DAVV Auditorium',
			coordinator: 'Prof. Anjali Sharma',
			status: 'Open'
		}
	];

	// Right sidebar registration deadlines data
	const deadlines = [
		{ name: 'Hackathon 2026', days: '3d left', color: 'bg-rose-50 text-rose-700 border-rose-100' },
		{
			name: 'Debate Championship',
			days: '5d left',
			color: 'bg-amber-50 text-amber-700 border-amber-100'
		},
		{
			name: 'Athletics Meet',
			days: '7d left',
			color: 'bg-emerald-50 text-emerald-700 border-emerald-100'
		},
		{
			name: 'Blood Donation Camp',
			days: '9d left',
			color: 'bg-blue-50 text-blue-700 border-blue-100'
		}
	];

	// Right sidebar categories counts
	const categoriesCounts = [
		{ name: 'Technical', count: 8, dotColor: 'bg-blue-600' },
		{ name: 'Research', count: 5, dotColor: 'bg-purple-600' },
		{ name: 'Sports', count: 6, dotColor: 'bg-emerald-600' },
		{ name: 'Cultural', count: 4, dotColor: 'bg-pink-600' },
		{ name: 'Leadership', count: 3, dotColor: 'bg-teal-600' },
		{ name: 'Social Service', count: 4, dotColor: 'bg-rose-600' },
		{ name: 'Public Speaking', count: 4, dotColor: 'bg-amber-600' }
	];

	// Enrollment interactive states
	let enrolledIds = $state<number[]>([]);

	// Base stats
	const baseEnrolledCount = 3;
	const basePendingCount = 1;
	const baseCompletedCount = 8;
	const baseCreditsEarned = 82;

	// Derived stats responsive to enrollments
	let enrolledCount = $derived(baseEnrolledCount + enrolledIds.length);
	let creditsEarned = $derived(
		baseCreditsEarned +
			enrolledIds.reduce((total, id) => {
				const act = allActivities.find((a) => a.id === id);
				return total + (act ? act.credits : 0);
			}, 0)
	);

	// Filters State
	let showFilters = $state(true);
	let tempFilters = $state({
		search: '',
		category: '',
		status: '',
		credits: '',
		mode: ''
	});
	let activeFilters = $state({
		search: '',
		category: '',
		status: '',
		credits: '',
		mode: ''
	});

	// Pagination State
	let currentPage = $state(1);
	const itemsPerPage = 6;

	// Filter logic
	let filteredActivities = $derived.by(() => {
		return allActivities.filter((activity) => {
			// 1. Search Query
			if (activeFilters.search) {
				const query = activeFilters.search.toLowerCase();
				const matchesName = activity.name.toLowerCase().includes(query);
				const matchesDesc = activity.description.toLowerCase().includes(query);
				if (!matchesName && !matchesDesc) return false;
			}

			// 2. Category
			if (activeFilters.category && activeFilters.category !== 'ALL') {
				if (activity.category !== activeFilters.category) return false;
			}

			// 3. Status
			if (activeFilters.status) {
				if (activity.status !== activeFilters.status) return false;
			}

			// 4. Credits Range
			if (activeFilters.credits) {
				if (activeFilters.credits === '1-10' && (activity.credits < 1 || activity.credits > 10))
					return false;
				if (activeFilters.credits === '11-15' && (activity.credits < 11 || activity.credits > 15))
					return false;
				if (activeFilters.credits === '16-20' && (activity.credits < 16 || activity.credits > 20))
					return false;
			}

			// 5. Mode
			if (activeFilters.mode) {
				if (activity.mode !== activeFilters.mode) return false;
			}

			return true;
		});
	});

	// Paginated activities
	let paginatedActivities = $derived.by(() => {
		const start = (currentPage - 1) * itemsPerPage;
		return filteredActivities.slice(start, start + itemsPerPage);
	});

	let totalPages = $derived(Math.ceil(filteredActivities.length / itemsPerPage) || 1);

	// Actions
	function handleApplyFilters() {
		activeFilters = { ...tempFilters };
		currentPage = 1;
	}

	function handleClearFilters() {
		tempFilters = { search: '', category: '', status: '', credits: '', mode: '' };
		activeFilters = { search: '', category: '', status: '', credits: '', mode: '' };
		currentPage = 1;
	}

	function handleRecommendedClick(category: string) {
		tempFilters.category = category.toUpperCase();
		handleApplyFilters();
	}

	function handleCategoryWidgetClick(category: string) {
		tempFilters.category = category.toUpperCase();
		handleApplyFilters();
	}

	function toggleEnrollment(id: number) {
		if (enrolledIds.includes(id)) {
			enrolledIds = enrolledIds.filter((itemId) => itemId !== id);
		} else {
			enrolledIds = [...enrolledIds, id];
		}
	}

	function getCategoryBadgeColor(category: string) {
		switch (category) {
			case 'TECHNICAL':
				return 'bg-slate-100 text-slate-700 border-slate-205 text-[10px]';
			case 'SPORTS':
				return 'bg-emerald-50 text-emerald-700 border-emerald-200 text-[10px]';
			case 'RESEARCH':
				return 'bg-sky-50 text-sky-700 border-sky-200 text-[10px]';
			case 'PUBLIC SPEAKING':
				return 'bg-amber-50 text-amber-700 border-amber-200 text-[10px]';
			case 'SOCIAL SERVICE':
				return 'bg-rose-50 text-rose-700 border-rose-200 text-[10px]';
			case 'CULTURAL':
				return 'bg-purple-50 text-purple-700 border-purple-200 text-[10px]';
			default:
				return 'bg-slate-100 text-slate-700 border-slate-200 text-[10px]';
		}
	}
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start font-sans">
	<!-- LEFT COLUMN (Main browse grid & search) -->
	<div class="lg:col-span-9 space-y-6">
		<!-- Find Activities Panel -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs">
			<div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-4">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Find Activities
				</h3>
				<button
					onclick={() => (showFilters = !showFilters)}
					class="text-xs font-bold text-[#881B1B] hover:underline"
				>
					{showFilters ? 'Hide Filters' : 'Show Filters'}
				</button>
			</div>

			{#if showFilters}
				<div transition:slide={{ duration: 200 }} class="space-y-4">
					<!-- Text search -->
					<div class="relative">
						<input
							type="text"
							bind:value={tempFilters.search}
							placeholder="Search activities, competitions, workshops..."
							class="w-full pl-3 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
						/>
					</div>

					<!-- Dropdowns grid -->
					<div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
						<!-- Category -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-category"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
								>Category</label
							>
							<select
								id="filter-category"
								bind:value={tempFilters.category}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Categories</option>
								<option value="TECHNICAL">Technical</option>
								<option value="RESEARCH">Research</option>
								<option value="SPORTS">Sports</option>
								<option value="CULTURAL">Cultural</option>
								<option value="SOCIAL SERVICE">Social Service</option>
								<option value="PUBLIC SPEAKING">Public Speaking</option>
								<option value="LEADERSHIP">Leadership</option>
							</select>
						</div>

						<!-- Status -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-status"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Status</label
							>
							<select
								id="filter-status"
								bind:value={tempFilters.status}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All</option>
								<option value="Open">Open</option>
								<option value="Closing Soon">Closing Soon</option>
								<option value="Closed">Closed</option>
							</select>
						</div>

						<!-- Credits Range -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-credits"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
								>Credits Range</label
							>
							<select
								id="filter-credits"
								bind:value={tempFilters.credits}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Credits</option>
								<option value="1-10">1 - 10 Points</option>
								<option value="11-15">11 - 15 Points</option>
								<option value="16-20">16 - 20 Points</option>
							</select>
						</div>

						<!-- Mode -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-mode"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Mode</label
							>
							<select
								id="filter-mode"
								bind:value={tempFilters.mode}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Modes</option>
								<option value="Offline">Offline</option>
								<option value="Online">Online</option>
								<option value="Hybrid">Hybrid</option>
							</select>
						</div>
					</div>

					<!-- Buttons -->
					<div class="flex items-center gap-3 pt-2">
						<button
							onclick={handleApplyFilters}
							class="px-5 py-2 bg-[#0B1535] hover:bg-[#1a2b5e] text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer"
						>
							Apply Filters
						</button>
						<button
							onclick={handleClearFilters}
							class="px-5 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-bold bg-white transition duration-200 cursor-pointer"
						>
							Clear Filters
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Recommended For You Row -->
		<div class="space-y-2.5">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
				Recommended For You <span class="text-[10px] text-slate-400 normal-case font-medium"
					>Based on your enrolled tracks</span
				>
			</h3>
			<div class="flex flex-wrap gap-2.5">
				<button
					onclick={() => handleRecommendedClick('TECHNICAL')}
					class="flex items-center gap-2 px-4 py-2 border border-[#881B1B]/20 bg-red-50 text-[#881B1B] rounded-lg text-xs font-bold hover:bg-[#881B1B]/10 transition cursor-pointer"
				>
					<!-- Sparkles Icon -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-4 h-4"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M9.813 15.904L9 21l-.813-5.096L3 15l5.096-.813L9 9l.813 5.096L15 15l-5.187.904zM18 5.25L17.25 8 16.5 5.25 13.75 4.5l2.75-.75L17.25 1l.75 2.75 2.75.75-2.75.75z"
						/>
					</svg>
					Technical Skill Track
					<span class="bg-[#881B1B] text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>6</span
					>
				</button>
				<button
					onclick={() => handleRecommendedClick('RESEARCH')}
					class="flex items-center gap-2 px-4 py-2 border border-purple-200 bg-purple-50 text-purple-700 rounded-lg text-xs font-bold hover:bg-purple-100 transition cursor-pointer"
				>
					<!-- Beaker/Science Icon -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-4 h-4"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M9.75 3.104v11.896m5.25-11.896v11.896M2.25 21h19.5M10 3.75h4"
						/>
					</svg>
					Research & Innovation
					<span class="bg-purple-700 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>4</span
					>
				</button>
				<button
					onclick={() => handleRecommendedClick('PUBLIC SPEAKING')}
					class="flex items-center gap-2 px-4 py-2 border border-amber-200 bg-amber-50 text-amber-700 rounded-lg text-xs font-bold hover:bg-amber-100 transition cursor-pointer"
				>
					<!-- Mic Icon -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-4 h-4"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M12 18.75a6 6 0 0 0 6-6v-1.5m-6 7.5a6 6 0 0 1-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 0 0 3-3V4.5a3 3 0 0 0-6 0v8.25a3 3 0 0 0 3 3z"
						/>
					</svg>
					Public Speaking
					<span class="bg-amber-700 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>3</span
					>
				</button>
			</div>
		</div>

		<!-- Activities Grid Section -->
		<div class="space-y-4">
			<div class="flex items-baseline justify-between border-b border-slate-200 pb-2">
				<h2 class="text-base font-bold font-serif text-[#0B1535]">
					All Activities <span class="text-xs font-semibold text-slate-400 normal-case ml-2"
						>{filteredActivities.length} found</span
					>
				</h2>
			</div>

			{#if paginatedActivities.length > 0}
				<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
					{#each paginatedActivities as activity (activity.id)}
						<div
							class="bg-white rounded-2xl border border-slate-200 overflow-hidden flex flex-col justify-between shadow-xs hover:shadow-md transition-all duration-200 relative group"
						>
							<!-- Top Banner / Color strip depending on category -->
							<div class="h-1.5 bg-[#881B1B]/80"></div>

							<div class="p-5 flex-grow space-y-4">
								<!-- Category Badge -->
								<div class="flex justify-between items-center">
									<span
										class="inline-flex px-2 py-0.5 rounded border text-[9px] font-bold uppercase tracking-wider {getCategoryBadgeColor(
											activity.category
										)}"
									>
										{activity.category}
									</span>

									{#if activity.status === 'Closing Soon'}
										<span
											class="text-[9px] font-bold text-rose-600 bg-rose-50 border border-rose-100 px-2 py-0.5 rounded-full uppercase tracking-wider"
										>
											Closing Soon
										</span>
									{/if}
								</div>

								<!-- Title & Description -->
								<div>
									<h4
										class="text-xs font-bold text-[#0B1535] leading-snug group-hover:text-[#881B1B] transition-colors"
									>
										{activity.name}
									</h4>
									<p class="text-[11px] text-slate-500 mt-2 leading-relaxed">
										{activity.description}
									</p>
								</div>

								<!-- Details List -->
								<div
									class="space-y-2 pt-3 border-t border-slate-100 text-[11px] font-medium text-slate-500"
								>
									<!-- Credits -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M11.48 3.499c.15-.358.682-.358.832 0l1.97 3.99a.75.75 0 0 0 .58.58l4.42 1.405c.4.082.56.58.268.88l-3.2 3.12a.75.75 0 0 0-.214.655l.8 4.417c.07.417-.37.74-.75.54L12 18.002l-3.95 2.06c-.38.2-.85-.122-.75-.54l.8-4.417a.75.75 0 0 0-.214-.655L4.69 11.235c-.29-.3-.13-.8.268-.88l4.42-1.405a.75.75 0 0 0 .58-.58l1.97-3.99z"
											/>
										</svg>
										<span
											>Credits: <strong class="text-slate-800">{activity.credits} Points</strong
											></span
										>
									</div>

									<!-- Mode -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M9.75 17L9 20l-1 3h4l-.5-3 .5-3m6-12V15a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 15V5.25m18 0A2.25 2.25 0 0 0 18.75 3H5.25A2.25 2.25 0 0 0 3 5.25"
											/>
										</svg>
										<span>Mode: <strong class="text-slate-800">{activity.mode}</strong></span>
									</div>

									<!-- Registration Deadline -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"
											/>
										</svg>
										<span
											>Registration Deadline: <strong class="text-slate-800"
												>{activity.regDeadline}</strong
											></span
										>
									</div>

									<!-- Activity Date -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
											/>
										</svg>
										<span
											>Activity Date: <strong class="text-slate-800">{activity.activityDate}</strong
											></span
										>
									</div>

									<!-- Venue -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
											/>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25s-7.5-4.108-7.5-11.25A7.5 7.5 0 1119.5 10.5z"
											/>
										</svg>
										<span>Venue: <strong class="text-slate-800">{activity.venue}</strong></span>
									</div>

									<!-- Coordinator -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0zM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632z"
											/>
										</svg>
										<span
											>Coordinator: <strong class="text-slate-800">{activity.coordinator}</strong
											></span
										>
									</div>
								</div>
							</div>

							<!-- Button -->
							<div class="px-5 pb-5">
								{#if enrolledIds.includes(activity.id)}
									<button
										onclick={() => toggleEnrollment(activity.id)}
										class="w-full py-2.5 bg-[#881B1B] hover:bg-[#731717] text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer flex items-center justify-center gap-1.5 shadow-sm font-sans uppercase tracking-wider"
									>
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
												d="M4.5 12.75l6 6 9-13.5"
											/>
										</svg>
										Enrolled
									</button>
								{:else}
									<button
										onclick={() => toggleEnrollment(activity.id)}
										disabled={activity.status === 'Closed'}
										class="w-full py-2.5 bg-[#881B1B] hover:bg-[#731717] disabled:bg-slate-200 disabled:text-slate-400 text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer shadow-xs focus:outline-none font-sans uppercase tracking-wider"
									>
										{activity.status === 'Closed' ? 'Registration Closed' : 'Enroll Now'}
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination Controls -->
				{#if totalPages > 1}
					<div
						class="flex items-center justify-center gap-2 pt-6 border-t border-slate-200 mt-6 font-sans"
					>
						<button
							disabled={currentPage === 1}
							onclick={() => (currentPage -= 1)}
							class="px-3 py-1.5 border border-slate-200 rounded-lg text-xs font-bold bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							&lt; Previous
						</button>

						{#each Array(totalPages) as _, index}
							<button
								onclick={() => (currentPage = index + 1)}
								class="w-8 h-8 rounded-lg text-xs font-bold transition flex items-center justify-center border {currentPage ===
								index + 1
									? 'bg-[#0B1535] border-[#0B1535] text-white'
									: 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'}"
							>
								{index + 1}
							</button>
						{/each}

						<button
							disabled={currentPage === totalPages}
							onclick={() => (currentPage += 1)}
							class="px-3 py-1.5 border border-slate-200 rounded-lg text-xs font-bold bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							Next &gt;
						</button>
					</div>
				{/if}
			{:else}
				<div class="bg-white border border-slate-200 rounded-xl p-12 text-center shadow-xs">
					<p class="text-slate-400 font-semibold text-xs">
						No activities found matching the filters.
					</p>
					<button
						onclick={handleClearFilters}
						class="text-[#881B1B] font-bold text-xs mt-3 hover:underline"
					>
						Clear Filters & Reset
					</button>
				</div>
			{/if}
		</div>
	</div>

	<!-- RIGHT COLUMN (Widgets sidebar) -->
	<div class="lg:col-span-3 space-y-6">
		<!-- Registration Deadlines Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Registration Deadlines
				</h3>
			</div>
			<div class="space-y-3 pt-1">
				{#each deadlines as item}
					<div class="flex items-center justify-between text-xs font-semibold text-slate-800">
						<span class="truncate pr-2 font-serif">{item.name}</span>
						<span
							class="px-2 py-0.5 rounded text-[10px] font-bold border shrink-0 font-sans {item.color}"
						>
							{item.days}
						</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- My Enrollment Status Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					My Enrollment Status
				</h3>
			</div>
			<div class="space-y-3 pt-1 text-xs font-semibold text-slate-500">
				<div class="flex justify-between items-center">
					<span>Enrolled</span>
					<span class="text-slate-900 font-bold">{enrolledCount}</span>
				</div>
				<div class="flex justify-between items-center">
					<span>Pending Approval</span>
					<span class="text-slate-900 font-bold">{basePendingCount}</span>
				</div>
				<div class="flex justify-between items-center">
					<span>Completed</span>
					<span class="text-slate-900 font-bold">{baseCompletedCount}</span>
				</div>
				<div class="flex justify-between items-center pt-2 border-t border-slate-100">
					<span>Credits Earned</span>
					<span class="text-[#881B1B] font-extrabold text-sm">{creditsEarned}</span>
				</div>
			</div>
		</div>

		<!-- Activity Categories Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Activity Categories
				</h3>
			</div>
			<div class="space-y-2.5 pt-1">
				{#each categoriesCounts as cat}
					<button
						onclick={() => handleCategoryWidgetClick(cat.name)}
						class="w-full flex items-center justify-between text-xs font-semibold text-slate-600 hover:text-slate-900 hover:bg-slate-50 p-1.5 rounded transition group cursor-pointer"
					>
						<div class="flex items-center gap-2">
							<span class="w-1.5 h-1.5 rounded-full {cat.dotColor}"></span>
							<span class="group-hover:font-bold">{cat.name}</span>
						</div>
						<span
							class="bg-slate-100 text-slate-700 font-bold text-[10px] px-2 py-0.5 rounded-full font-sans"
						>
							{cat.count}
						</span>
					</button>
				{/each}
			</div>
		</div>
	</div>
</div>
