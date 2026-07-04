<script lang="ts">
	// Define interfaces
	interface Student {
		rank: string;
		initials: string;
		name: string;
		course: string;
		sem: number;
		activities: number;
		credits: number;
		grade: 'O' | 'A' | 'B';
		isSelf?: boolean;
		avatarBg: string;
	}

	interface Champion {
		track: string;
		name: string;
		credits: number;
		initials: string;
		avatarBg: string;
	}

	interface RecognitionLevel {
		name: string;
		creditsRequired: number;
		percentile: string;
		icon: 'star' | 'check';
		theme: 'gold' | 'rose' | 'blue' | 'slate';
		unlocked: boolean;
	}

	interface Badge {
		name: string;
		description: string;
		date?: string;
		icon:
			| 'bolt'
			| 'star'
			| 'shield'
			| 'consistent'
			| 'microphone'
			| 'book'
			| 'music'
			| 'trophy'
			| 'heart'
			| 'award';
		theme: 'gold' | 'navy' | 'green' | 'purple' | 'blue' | 'pink' | 'teal' | 'slate';
		locked: boolean;
	}

	// Filter state for Academic Year
	let selectedYear = $state('2025-26');

	// Mock Student Data (changes slightly based on Year to show interactivity)
	const students2025: Student[] = [
		{
			rank: '01',
			initials: 'AS',
			name: 'Aarav Sharma',
			course: 'Int. MCA',
			sem: 6,
			activities: 28,
			credits: 145,
			grade: 'O',
			avatarBg: 'bg-amber-100 text-amber-800 border-amber-300'
		},
		{
			rank: '02',
			initials: 'PP',
			name: 'Priya Patel',
			course: 'Int. MCA',
			sem: 6,
			activities: 26,
			credits: 138,
			grade: 'O',
			avatarBg: 'bg-purple-100 text-purple-800 border-purple-300'
		},
		{
			rank: '03',
			initials: 'DR',
			name: 'Deepak Rathore',
			course: 'Int. MCA',
			sem: 5,
			activities: 24,
			credits: 130,
			grade: 'A',
			avatarBg: 'bg-orange-100 text-orange-800 border-orange-300'
		},
		{
			rank: '04',
			initials: 'RV',
			name: 'Rahul Verma',
			course: 'B.Tech CSE',
			sem: 6,
			activities: 24,
			credits: 118,
			grade: 'A',
			isSelf: true,
			avatarBg: 'bg-red-100 text-red-800 border-red-300'
		},
		{
			rank: '05',
			initials: 'SK',
			name: 'Sneha Kulkarni',
			course: 'B.Tech CSE',
			sem: 6,
			activities: 22,
			credits: 112,
			grade: 'A',
			avatarBg: 'bg-teal-100 text-teal-800 border-teal-300'
		},
		{
			rank: '06',
			initials: 'AM',
			name: 'Arjun Mehta',
			course: 'MCA',
			sem: 4,
			activities: 20,
			credits: 108,
			grade: 'A',
			avatarBg: 'bg-blue-100 text-blue-800 border-blue-300'
		},
		{
			rank: '07',
			initials: 'PD',
			name: 'Pooja Desai',
			course: 'MBA (MS)',
			sem: 4,
			activities: 19,
			credits: 102,
			grade: 'A',
			avatarBg: 'bg-pink-100 text-pink-800 border-pink-300'
		},
		{
			rank: '08',
			initials: 'RJ',
			name: 'Rohit Jaiswal',
			course: 'B.Tech CSE',
			sem: 5,
			activities: 18,
			credits: 96,
			grade: 'B',
			avatarBg: 'bg-indigo-100 text-indigo-800 border-indigo-300'
		},
		{
			rank: '09',
			initials: 'KN',
			name: 'Kavya Nair',
			course: 'Int. MCA',
			sem: 4,
			activities: 17,
			credits: 88,
			grade: 'B',
			avatarBg: 'bg-fuchsia-100 text-fuchsia-800 border-fuchsia-300'
		},
		{
			rank: '10',
			initials: 'MS',
			name: 'Manish Soni',
			course: 'MCA',
			sem: 6,
			activities: 16,
			credits: 82,
			grade: 'B',
			avatarBg: 'bg-emerald-100 text-emerald-800 border-emerald-300'
		},
		{
			rank: '11',
			initials: 'NT',
			name: 'Nisha Trivedi',
			course: 'MBA',
			sem: 4,
			activities: 14,
			credits: 76,
			grade: 'B',
			avatarBg: 'bg-sky-100 text-sky-800 border-sky-300'
		},
		{
			rank: '12',
			initials: 'SG',
			name: 'Sachin Gupta',
			course: 'B.Tech IT',
			sem: 5,
			activities: 13,
			credits: 71,
			grade: 'B',
			avatarBg: 'bg-slate-105 text-slate-800 border-slate-300'
		}
	];

	const students2024: Student[] = [
		{
			rank: '01',
			initials: 'PP',
			name: 'Priya Patel',
			course: 'Int. MCA',
			sem: 4,
			activities: 24,
			credits: 128,
			grade: 'A',
			avatarBg: 'bg-purple-100 text-purple-800 border-purple-300'
		},
		{
			rank: '02',
			initials: 'AS',
			name: 'Aarav Sharma',
			course: 'Int. MCA',
			sem: 4,
			activities: 22,
			credits: 122,
			grade: 'A',
			avatarBg: 'bg-amber-100 text-amber-800 border-amber-300'
		},
		{
			rank: '03',
			initials: 'SK',
			name: 'Sneha Kulkarni',
			course: 'B.Tech CSE',
			sem: 4,
			activities: 21,
			credits: 115,
			grade: 'A',
			avatarBg: 'bg-teal-100 text-teal-800 border-teal-300'
		},
		{
			rank: '04',
			initials: 'RV',
			name: 'Rahul Verma',
			course: 'B.Tech CSE',
			sem: 4,
			activities: 18,
			credits: 104,
			grade: 'A',
			isSelf: true,
			avatarBg: 'bg-red-100 text-red-800 border-red-300'
		},
		{
			rank: '05',
			initials: 'DR',
			name: 'Deepak Rathore',
			course: 'Int. MCA',
			sem: 3,
			activities: 19,
			credits: 101,
			grade: 'A',
			avatarBg: 'bg-orange-100 text-orange-800 border-orange-300'
		},
		{
			rank: '06',
			initials: 'PD',
			name: 'Pooja Desai',
			course: 'MBA (MS)',
			sem: 2,
			activities: 16,
			credits: 92,
			grade: 'B',
			avatarBg: 'bg-pink-100 text-pink-800 border-pink-300'
		},
		{
			rank: '07',
			initials: 'RJ',
			name: 'Rohit Jaiswal',
			course: 'B.Tech CSE',
			sem: 3,
			activities: 15,
			credits: 89,
			grade: 'B',
			avatarBg: 'bg-indigo-100 text-indigo-800 border-indigo-300'
		},
		{
			rank: '08',
			initials: 'AM',
			name: 'Arjun Mehta',
			course: 'MCA',
			sem: 2,
			activities: 14,
			credits: 85,
			grade: 'B',
			avatarBg: 'bg-blue-100 text-blue-800 border-blue-300'
		},
		{
			rank: '09',
			initials: 'KN',
			name: 'Kavya Nair',
			course: 'Int. MCA',
			sem: 2,
			activities: 12,
			credits: 78,
			grade: 'B',
			avatarBg: 'bg-fuchsia-100 text-fuchsia-800 border-fuchsia-300'
		},
		{
			rank: '10',
			initials: 'MS',
			name: 'Manish Soni',
			course: 'MCA',
			sem: 4,
			activities: 11,
			credits: 72,
			grade: 'B',
			avatarBg: 'bg-emerald-100 text-emerald-800 border-emerald-300'
		},
		{
			rank: '11',
			initials: 'NT',
			name: 'Nisha Trivedi',
			course: 'MBA',
			sem: 2,
			activities: 9,
			credits: 65,
			grade: 'B',
			avatarBg: 'bg-sky-100 text-sky-800 border-sky-300'
		},
		{
			rank: '12',
			initials: 'SG',
			name: 'Sachin Gupta',
			course: 'B.Tech IT',
			sem: 3,
			activities: 8,
			credits: 60,
			grade: 'B',
			avatarBg: 'bg-slate-105 text-slate-800 border-slate-300'
		}
	];

	// Derived state for the active list
	let activeStudents = $derived(selectedYear === '2025-26' ? students2025 : students2024);

	// Derived Podium Students
	let podiumFirst = $derived(activeStudents.find((s) => s.rank === '01')!);
	let podiumSecond = $derived(activeStudents.find((s) => s.rank === '02')!);
	let podiumThird = $derived(activeStudents.find((s) => s.rank === '03')!);

	// Derived Rahul Verma (YOU) credits to show dynamic Recognition Levels
	let currentUserCredits = $derived(activeStudents.find((s) => s.isSelf)?.credits || 118);

	// Category Champions Data
	const champions2025: Champion[] = [
		{
			track: 'TECHNICAL',
			name: 'Aarav Sharma',
			credits: 52,
			initials: 'AS',
			avatarBg: 'bg-blue-50 text-blue-700 border-blue-200'
		},
		{
			track: 'PUBLIC SPEAKING',
			name: 'Priya Patel',
			credits: 36,
			initials: 'PP',
			avatarBg: 'bg-purple-50 text-purple-700 border-purple-200'
		},
		{
			track: 'RESEARCH',
			name: 'Deepak Rathore',
			credits: 28,
			initials: 'DR',
			avatarBg: 'bg-orange-50 text-orange-700 border-orange-200'
		},
		{
			track: 'SPORTS',
			name: 'Rohit Jaiswal',
			credits: 24,
			initials: 'RJ',
			avatarBg: 'bg-teal-50 text-teal-700 border-teal-200'
		},
		{
			track: 'SOCIAL SERVICE',
			name: 'Pooja Desai',
			credits: 22,
			initials: 'PD',
			avatarBg: 'bg-rose-50 text-rose-700 border-rose-200'
		}
	];

	const champions2024: Champion[] = [
		{
			track: 'TECHNICAL',
			name: 'Priya Patel',
			credits: 42,
			initials: 'PP',
			avatarBg: 'bg-purple-50 text-purple-700 border-purple-200'
		},
		{
			track: 'PUBLIC SPEAKING',
			name: 'Aarav Sharma',
			credits: 32,
			initials: 'AS',
			avatarBg: 'bg-blue-50 text-blue-700 border-blue-200'
		},
		{
			track: 'RESEARCH',
			name: 'Sneha Kulkarni',
			credits: 25,
			initials: 'SK',
			avatarBg: 'bg-teal-50 text-teal-700 border-teal-200'
		},
		{
			track: 'SPORTS',
			name: 'Rahul Verma',
			credits: 24,
			initials: 'RV',
			avatarBg: 'bg-red-50 text-red-700 border-red-200'
		},
		{
			track: 'SOCIAL SERVICE',
			name: 'Pooja Desai',
			credits: 18,
			initials: 'PD',
			avatarBg: 'bg-rose-50 text-rose-700 border-rose-200'
		}
	];

	let activeChampions = $derived(selectedYear === '2025-26' ? champions2025 : champions2024);

	// Recognition Levels calculation based on current user credits
	let recognitionLevels = $derived<RecognitionLevel[]>([
		{
			name: 'Top Performer',
			creditsRequired: 140,
			percentile: 'Top 1%',
			icon: 'star',
			theme: 'gold',
			unlocked: currentUserCredits >= 140
		},
		{
			name: 'Outstanding Contributor',
			creditsRequired: 100,
			percentile: 'Top 5%',
			icon: 'check',
			theme: 'rose',
			unlocked: currentUserCredits >= 100
		},
		{
			name: 'Academic Achiever',
			creditsRequired: 80,
			percentile: 'Top 10%',
			icon: 'check',
			theme: 'blue',
			unlocked: currentUserCredits >= 80
		},
		{
			name: 'Active Participant',
			creditsRequired: 50,
			percentile: 'Top 25%',
			icon: 'check',
			theme: 'slate',
			unlocked: currentUserCredits >= 50
		}
	]);

	// Achievement Badges
	const badges: Badge[] = [
		{
			name: 'Hackathon Hero',
			description: 'Won a national-level hackathon',
			date: '12 Jan 2026',
			icon: 'bolt',
			theme: 'gold',
			locked: false
		},
		{
			name: 'Top 5 Achiever',
			description: 'Ranked within Top 5 students',
			date: '18 Jun 2026',
			icon: 'star',
			theme: 'gold',
			locked: false
		},
		{
			name: 'Century Club',
			description: 'Crossed 100 verified credits',
			date: '03 May 2026',
			icon: 'shield',
			theme: 'navy',
			locked: false
		},
		{
			name: 'Consistent Performer',
			description: 'Earned credits every semester',
			date: 'Jan 2026',
			icon: 'consistent',
			theme: 'green',
			locked: false
		},
		{
			name: 'Silver Orator',
			description: '24+ Public Speaking credits',
			date: '02 May 2025',
			icon: 'microphone',
			theme: 'purple',
			locked: false
		},
		{
			name: 'Research Scholar',
			description: 'Presented at 3+ symposia',
			date: '20 Jan 2025',
			icon: 'book',
			theme: 'blue',
			locked: false
		},
		{
			name: 'Cultural Catalyst',
			description: 'Participated in 3+ cultural events',
			date: '10 May 2025',
			icon: 'music',
			theme: 'pink',
			locked: false
		},
		{
			name: 'Sports Star',
			description: 'Won a position in sports meet',
			date: '08 Jul 2025',
			icon: 'trophy',
			theme: 'teal',
			locked: false
		},
		{
			name: 'Community Pillar',
			description: 'Completed 3 social service drives',
			icon: 'heart',
			theme: 'slate',
			locked: true
		},
		{
			name: 'Grade O Achiever',
			description: 'Reach 140+ verified credits',
			icon: 'award',
			theme: 'slate',
			locked: true
		}
	];

	function getGradeColors(grade: 'O' | 'A' | 'B') {
		switch (grade) {
			case 'O':
				return {
					badge: 'bg-emerald-50 text-emerald-700 border-emerald-200',
					underline: 'border-emerald-500'
				};
			case 'A':
				return {
					badge: 'bg-rose-50 text-rose-700 border-rose-200',
					underline: 'border-rose-500'
				};
			case 'B':
				return {
					badge: 'bg-blue-50 text-blue-700 border-blue-200',
					underline: 'border-blue-500'
				};
		}
	}
</script>

<div class="space-y-6 font-sans">
	<!-- ==================== 1. TOP PERFORMERS PODIUM ==================== -->
	<section
		class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs relative"
		aria-label="Top Performers Podium"
	>
		<div class="flex items-center justify-between pb-6 border-b border-slate-100 mb-6">
			<div>
				<h2 class="text-base font-bold font-serif text-[#0B1535]">Top Performers</h2>
				<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">
					Ranked by verified extracurricular credits
				</p>
			</div>

			<!-- Year Dropdown -->
			<div class="relative">
				<select
					bind:value={selectedYear}
					class="appearance-none pl-3 pr-8 py-1.5 bg-slate-50 border border-slate-205 rounded-lg text-xs font-bold text-slate-700 focus:outline-none focus:border-[#881B1B] cursor-pointer transition-colors"
				>
					<option value="2025-26">2025-26</option>
					<option value="2024-25">2024-25</option>
				</select>
				<span
					class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-500 pointer-events-none"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2.5"
						stroke="currentColor"
						class="w-3.5 h-3.5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
					</svg>
				</span>
			</div>
		</div>

		<!-- Podium Layout Grid -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 items-end max-w-3xl mx-auto pt-6 pb-2">
			<!-- 2nd Place (Left) -->
			{#if podiumSecond}
				<div class="flex flex-col items-center order-2 md:order-1 mt-4 md:mt-0">
					<!-- Circular Rank Badge -->
					<div
						class="w-7 h-7 rounded-full bg-slate-100 border border-slate-250 flex items-center justify-center text-xs font-bold text-slate-500 mb-3 shadow-3xs"
					>
						II
					</div>
					<!-- Card Container -->
					<div
						class="w-full bg-white border border-slate-200 rounded-xl p-5 shadow-xs flex flex-col items-center text-center relative hover:shadow-md transition-shadow duration-200"
					>
						<!-- Initials Avatar -->
						<div
							class="w-14 h-14 rounded-full bg-purple-50 text-purple-700 border-2 border-purple-200 flex items-center justify-center font-bold text-base shadow-sm shrink-0 mb-3"
						>
							{podiumSecond.initials}
						</div>
						<h3 class="text-xs font-bold text-[#0B1535] leading-tight">{podiumSecond.name}</h3>
						<p class="text-[9px] font-semibold text-slate-400 mt-1 uppercase tracking-wider">
							{podiumSecond.course} - Sem {podiumSecond.sem}
						</p>

						<div class="mt-3 flex items-baseline gap-1">
							<span class="text-xl font-extrabold text-[#0B1535] leading-none"
								>{podiumSecond.credits}</span
							>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest"
								>credits</span
							>
						</div>

						<span
							class="mt-2.5 px-2.5 py-0.5 rounded text-[9px] font-extrabold uppercase tracking-wide border {getGradeColors(
								podiumSecond.grade
							).badge}"
						>
							Grade {podiumSecond.grade}
						</span>
					</div>
					<!-- Podium block -->
					<div
						class="w-full h-10 bg-slate-100 border border-slate-200 border-t-0 rounded-b-xl flex items-center justify-center text-slate-400 font-extrabold text-sm shadow-3xs"
					>
						II
					</div>
				</div>
			{/if}

			<!-- 1st Place (Center) -->
			{#if podiumFirst}
				<div class="flex flex-col items-center order-1 md:order-2 scale-105 z-10">
					<!-- Circular Rank Badge -->
					<div
						class="w-7 h-7 rounded-full bg-amber-50 border border-amber-250 flex items-center justify-center text-xs font-bold text-amber-700 mb-3 shadow-3xs"
					>
						I
					</div>
					<!-- Card Container -->
					<div
						class="w-full bg-white border-2 border-amber-300 rounded-xl p-6 shadow-md flex flex-col items-center text-center relative hover:shadow-lg transition-shadow duration-200"
					>
						<!-- Gold crown or highlights decorator -->
						<div
							class="absolute -top-3.5 bg-amber-400 text-white rounded-full p-1 shadow-md border border-white"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="currentColor"
								viewBox="0 0 24 24"
								class="w-3.5 h-3.5"
							>
								<path
									d="M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3c1.606 0 3.012.723 3.966 1.892L12 5.344l.346-.452C13.298 3.723 14.704 3 16.312 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z"
								/>
							</svg>
						</div>
						<!-- Initials Avatar -->
						<div
							class="w-16 h-16 rounded-full bg-amber-50 text-amber-600 border-2 border-amber-300 flex items-center justify-center font-bold text-lg shadow-sm shrink-0 mb-3"
						>
							{podiumFirst.initials}
						</div>
						<h3 class="text-sm font-bold text-[#0B1535] leading-tight">{podiumFirst.name}</h3>
						<p class="text-[9px] font-semibold text-slate-400 mt-1 uppercase tracking-wider">
							{podiumFirst.course} - Sem {podiumFirst.sem}
						</p>

						<div class="mt-3 flex items-baseline gap-1">
							<span class="text-2xl font-extrabold text-[#0B1535] leading-none"
								>{podiumFirst.credits}</span
							>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest"
								>credits</span
							>
						</div>

						<span
							class="mt-2.5 px-2.5 py-0.5 rounded text-[9px] font-extrabold uppercase tracking-wide border {getGradeColors(
								podiumFirst.grade
							).badge}"
						>
							Grade {podiumFirst.grade}
						</span>
					</div>
					<!-- Podium block -->
					<div
						class="w-full h-14 bg-amber-50/70 border border-amber-200 border-t-0 rounded-b-xl flex items-center justify-center text-amber-600 font-extrabold text-sm shadow-3xs"
					>
						I
					</div>
				</div>
			{/if}

			<!-- 3rd Place (Right) -->
			{#if podiumThird}
				<div class="flex flex-col items-center order-3 md:order-3 mt-4 md:mt-0">
					<!-- Circular Rank Badge -->
					<div
						class="w-7 h-7 rounded-full bg-orange-50 border border-orange-255 flex items-center justify-center text-xs font-bold text-orange-850 mb-3 shadow-3xs"
					>
						III
					</div>
					<!-- Card Container -->
					<div
						class="w-full bg-white border border-slate-200 rounded-xl p-5 shadow-xs flex flex-col items-center text-center relative hover:shadow-md transition-shadow duration-200"
					>
						<!-- Initials Avatar -->
						<div
							class="w-14 h-14 rounded-full bg-orange-50 text-orange-700 border-2 border-orange-200 flex items-center justify-center font-bold text-base shadow-sm shrink-0 mb-3"
						>
							{podiumThird.initials}
						</div>
						<h3 class="text-xs font-bold text-[#0B1535] leading-tight">{podiumThird.name}</h3>
						<p class="text-[9px] font-semibold text-slate-400 mt-1 uppercase tracking-wider">
							{podiumThird.course} - Sem {podiumThird.sem}
						</p>

						<div class="mt-3 flex items-baseline gap-1">
							<span class="text-xl font-extrabold text-[#0B1535] leading-none"
								>{podiumThird.credits}</span
							>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest"
								>credits</span
							>
						</div>

						<span
							class="mt-2.5 px-2.5 py-0.5 rounded text-[9px] font-extrabold uppercase tracking-wide border {getGradeColors(
								podiumThird.grade
							).badge}"
						>
							Grade {podiumThird.grade}
						</span>
					</div>
					<!-- Podium block -->
					<div
						class="w-full h-8 bg-orange-50/50 border border-orange-200 border-t-0 rounded-b-xl flex items-center justify-center text-orange-700/80 font-extrabold text-sm shadow-3xs"
					>
						III
					</div>
				</div>
			{/if}
		</div>
	</section>

	<!-- ==================== 2. STUDENT RANKINGS TABLE ==================== -->
	<section
		class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs"
		aria-label="Student Rankings"
	>
		<div class="pb-4 border-b border-slate-105 mb-4">
			<h2 class="text-base font-bold font-serif text-[#0B1535]">Student Rankings</h2>
			<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">
				Showing {activeStudents.length} students - Sorted by verified credits
			</p>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr
						class="text-[10px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100 bg-slate-50/50"
					>
						<th class="py-3 px-4 w-16 text-center">Rank</th>
						<th class="py-3 px-4">Student</th>
						<th class="py-3 px-4">Course</th>
						<th class="py-3 px-4 text-center">Sem</th>
						<th class="py-3 px-4 text-center">Activities</th>
						<th class="py-3 px-4 text-center">Credits</th>
						<th class="py-3 px-4 text-center">Grade</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 font-sans">
					{#each activeStudents as student}
						<tr
							class="hover:bg-slate-50/75 transition-colors {student.isSelf
								? 'bg-red-50/30 border-y border-red-150/40 font-bold'
								: ''}"
						>
							<!-- Rank Column -->
							<td class="py-3 px-4 text-center">
								{#if student.rank === '01'}
									<span
										class="inline-flex w-7 h-7 rounded-full items-center justify-center bg-[#FDF9EA] text-[#C89B3C] border border-[#F5E6BE] text-[10px] font-extrabold shadow-3xs"
									>
										01
									</span>
								{:else if student.rank === '02'}
									<span
										class="inline-flex w-7 h-7 rounded-full items-center justify-center bg-[#F3F4F6] text-[#6B7280] border border-[#E5E7EB] text-[10px] font-extrabold shadow-3xs"
									>
										02
									</span>
								{:else if student.rank === '03'}
									<span
										class="inline-flex w-7 h-7 rounded-full items-center justify-center bg-[#FDF5EB] text-[#C0703C] border border-[#F5DCBE] text-[10px] font-extrabold shadow-3xs"
									>
										03
									</span>
								{:else}
									<span
										class="inline-flex w-7 h-7 rounded-full items-center justify-center border border-slate-200 text-slate-500 bg-slate-50/50 text-[10px] font-semibold"
									>
										{student.rank}
									</span>
								{/if}
							</td>

							<!-- Student Profile -->
							<td class="py-3 px-4">
								<div class="flex items-center gap-3">
									<div
										class="w-8 h-8 rounded-full border flex items-center justify-center text-[11px] font-bold shrink-0 shadow-3xs {student.avatarBg}"
									>
										{student.initials}
									</div>
									<div class="flex flex-col">
										<div class="flex items-center gap-1.5">
											<span class="text-slate-900 font-bold">{student.name}</span>
											{#if student.isSelf}
												<span
													class="inline-flex px-1.5 py-0.5 bg-red-100 text-red-700 text-[8px] font-extrabold uppercase rounded-sm tracking-wider"
												>
													YOU
												</span>
											{/if}
										</div>
									</div>
								</div>
							</td>

							<!-- Course -->
							<td class="py-3 px-4 text-slate-600 font-semibold">{student.course}</td>

							<!-- Sem -->
							<td class="py-3 px-4 text-center text-[#0b1535] font-bold">{student.sem}</td>

							<!-- Activities -->
							<td class="py-3 px-4 text-center text-slate-500 font-bold">{student.activities}</td>

							<!-- Credits (underlined with grade color) -->
							<td class="py-3 px-4 text-center">
								<span
									class="border-b-[3px] pb-0.5 font-extrabold text-[#0B1535] {getGradeColors(
										student.grade
									).underline}"
								>
									{student.credits}
								</span>
							</td>

							<!-- Grade -->
							<td class="py-3 px-4 text-center">
								<span
									class="inline-flex px-2 py-0.5 rounded text-[9px] font-extrabold border uppercase tracking-wider {getGradeColors(
										student.grade
									).badge}"
								>
									{student.grade}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</section>

	<!-- ==================== 3. CATEGORY CHAMPIONS ==================== -->
	<section class="space-y-3" aria-label="Category Champions">
		<div>
			<h2 class="text-base font-bold font-serif text-[#0B1535]">Category Champions</h2>
			<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">
				Top performer per track
			</p>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
			{#each activeChampions as champ}
				<div
					class="bg-white border border-slate-205 rounded-xl p-4 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow duration-200 relative overflow-hidden group"
				>
					<div
						class="absolute top-0 left-0 w-full h-1 bg-slate-100 group-hover:bg-[#881B1B] transition-colors"
					></div>
					<div class="space-y-3">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">
							{champ.track}
						</span>

						<div class="flex items-center gap-3">
							<!-- Initials Avatar -->
							<div
								class="w-9 h-9 rounded-full border flex items-center justify-center text-xs font-bold shrink-0 shadow-3xs {champ.avatarBg}"
							>
								{champ.initials}
							</div>
							<div>
								<h3 class="text-xs font-bold text-slate-800 leading-tight">{champ.name}</h3>
								<p class="text-[10px] font-semibold text-slate-500 mt-0.5">
									{champ.credits}
									<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide"
										>credits</span
									>
								</p>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- ==================== 4. RECOGNITION & BADGES ROW ==================== -->
	<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
		<!-- Recognition Levels (Left Column) -->
		<section
			class="lg:col-span-7 bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4"
			aria-label="Recognition Levels"
		>
			<div class="pb-3 border-b border-slate-100">
				<h2 class="text-base font-bold font-serif text-[#0B1535]">Recognition Levels</h2>
				<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">
					Levels unlocked based on credits
				</p>
			</div>

			<div class="space-y-3">
				{#each recognitionLevels as level}
					<div
						class="flex items-center justify-between p-4 border rounded-xl transition-all duration-200 gap-4
						{level.unlocked
							? level.theme === 'rose'
								? 'bg-rose-50/20 border-rose-200'
								: level.theme === 'blue'
									? 'bg-blue-50/20 border-blue-200'
									: level.theme === 'slate'
										? 'bg-slate-50/50 border-slate-200'
										: 'bg-white border-slate-200'
							: 'border-slate-150 bg-slate-50/20 opacity-60'}"
					>
						<div class="flex items-center gap-3">
							<!-- Indicator Icon Circle -->
							<div
								class="w-8 h-8 rounded-full flex items-center justify-center border shrink-0 shadow-3xs
								{level.unlocked
									? level.theme === 'rose'
										? 'bg-rose-50 border-rose-200 text-rose-600'
										: level.theme === 'blue'
											? 'bg-blue-50 border-blue-200 text-blue-600'
											: level.theme === 'slate'
												? 'bg-slate-105 border-slate-250 text-slate-600'
												: 'bg-white border-slate-200 text-slate-400'
									: 'bg-slate-100 border-slate-200 text-slate-400'}"
							>
								{#if level.icon === 'star'}
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
											d="M12 4.5v15m7.5-7.5h-15"
										/>
									</svg>
								{:else if level.icon === 'check'}
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
								{/if}
							</div>

							<div>
								<h3 class="text-xs font-bold text-slate-800 leading-tight">{level.name}</h3>
								<p class="text-[10px] text-slate-450 mt-1">
									{level.creditsRequired}+ credits required
								</p>
							</div>
						</div>

						<span
							class="text-[9px] font-extrabold uppercase px-2.5 py-1 rounded tracking-wide border shadow-3xs
							{level.unlocked
								? level.theme === 'rose'
									? 'bg-rose-100 border-rose-300 text-rose-700'
									: level.theme === 'blue'
										? 'bg-blue-100 border-blue-300 text-blue-700'
										: level.theme === 'slate'
											? 'bg-slate-100 border-slate-300 text-slate-700'
											: 'bg-slate-100 border-slate-200 text-slate-600'
								: 'bg-slate-50 border-slate-150 text-slate-400'}"
						>
							{level.percentile}
						</span>
					</div>
				{/each}
			</div>
		</section>

		<!-- Achievement Badges (Right Column) -->
		<section
			class="lg:col-span-5 bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4"
			aria-label="Achievement Badges"
		>
			<div class="pb-3 border-b border-slate-100 flex items-center justify-between">
				<div>
					<h2 class="text-base font-bold font-serif text-[#0B1535]">Achievement Badges</h2>
					<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">
						Earned through verified participation
					</p>
				</div>
				<!-- Progress Tracker Pill -->
				<span
					class="bg-[#0b1535] text-white text-[10px] font-bold px-2 py-0.5 rounded-full shadow-3xs"
				>
					8/10
				</span>
			</div>

			<!-- Scrollable Badges Container -->
			<div class="space-y-3.5 max-h-[432px] overflow-y-auto pr-1">
				<!-- Earned Section -->
				<div class="text-[9px] font-extrabold text-slate-400 uppercase tracking-wider block">
					Earned &middot; 8
				</div>

				{#each badges.filter((b) => !b.locked) as badge}
					<div
						class="flex items-center justify-between p-3 border border-slate-150 bg-slate-50/20 hover:bg-slate-50/50 transition-colors rounded-xl gap-4 relative overflow-hidden"
					>
						<div class="flex items-center gap-3">
							<!-- Colored Badge Circle -->
							<div
								class="w-9 h-9 rounded-full border flex items-center justify-center shrink-0 shadow-3xs
								{badge.theme === 'gold'
									? 'bg-amber-50 text-amber-600 border-amber-100'
									: badge.theme === 'navy'
										? 'bg-blue-50 text-blue-800 border-blue-200'
										: badge.theme === 'green'
											? 'bg-emerald-50 text-emerald-600 border-emerald-100'
											: badge.theme === 'purple'
												? 'bg-purple-50 text-purple-600 border-purple-150'
												: badge.theme === 'blue'
													? 'bg-sky-50 text-sky-600 border-sky-150'
													: badge.theme === 'pink'
														? 'bg-pink-50 text-pink-600 border-pink-150'
														: badge.theme === 'teal'
															? 'bg-teal-50 text-teal-600 border-teal-150'
															: 'bg-slate-50 text-slate-600 border-slate-200'}"
							>
								{#if badge.icon === 'bolt'}
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
											d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"
										/>
									</svg>
								{:else if badge.icon === 'star'}
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
											d="M11.48 3.499c.15-.358.682-.358.832 0l1.97 3.99a.75.75 0 0 0 .58.58l4.42 1.405c.4.082.56.58.268.88l-3.2 3.12a.75.75 0 0 0-.214.655l.8 4.417c.07.417-.37.74-.75.54L12 18.002l-3.95 2.06c-.38.2-.85-.122-.75-.54l.8-4.417a.75.75 0 0 0-.214-.655L4.69 11.235c-.29-.3-.13-.8.268-.88l4.42-1.405a.75.75 0 0 0 .58-.58l1.97-3.99z"
										/>
									</svg>
								{:else if badge.icon === 'shield'}
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
											d="M9 12.75L11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 0 1-1.043 3.296 3.745 3.745 0 0 1-3.296 1.043A3.745 3.745 0 0 1 12 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 0 1-3.296-1.043 3.745 3.745 0 0 1-1.043-3.296A3.745 3.745 0 0 1 3 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 0 1 1.043-3.296 3.745 3.745 0 0 1 3.296-1.043A3.746 3.746 0 0 1 12 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 0 1 3.296 1.043 3.746 3.746 0 0 1 1.043 3.296A3.745 3.745 0 0 1 21 12Z"
										/>
									</svg>
								{:else if badge.icon === 'consistent'}
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
											d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"
										/>
									</svg>
								{:else if badge.icon === 'microphone'}
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
								{:else if badge.icon === 'book'}
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
											d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
										/>
									</svg>
								{:else if badge.icon === 'music'}
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
											d="M9 9l10.5-3m0 0v11.25m0-11.25L9 9m10.5 5.25v3.75m0-3.75a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0zm-10.5 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0z"
										/>
									</svg>
								{:else if badge.icon === 'trophy'}
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
											d="M16.5 18.75h-9m9 0a3 3 0 0 1 3-3h.008a3 3 0 0 1 3 3m-15 0a3 3 0 0 0-3-3H4.5a3 3 0 0 0-3 3m15 0v-1.125c0-1.657-1.343-3-3-3h-4.5c-1.657 0-3 1.343-3 3v1.125m9-10.5a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"
										/>
									</svg>
								{/if}
							</div>

							<div>
								<h3 class="text-xs font-bold text-slate-805 leading-tight">{badge.name}</h3>
								<p class="text-[10px] text-slate-400 mt-0.5">{badge.description}</p>
							</div>
						</div>

						<div class="flex items-center gap-1.5 shrink-0">
							<span class="text-[9px] text-slate-400 whitespace-nowrap font-medium"
								>{badge.date}</span
							>
							<span
								class="w-4 h-4 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center shrink-0 border border-emerald-200"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="3"
									stroke="currentColor"
									class="w-2.5 h-2.5"
								>
									<path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
								</svg>
							</span>
						</div>
					</div>
				{/each}

				<!-- Locked Section -->
				<div class="text-[9px] font-extrabold text-slate-400 uppercase tracking-wider block pt-2">
					Locked &middot; 2
				</div>

				{#each badges.filter((b) => b.locked) as badge}
					<div
						class="flex items-center justify-between p-3 border border-dashed border-slate-200 bg-slate-50/10 opacity-50 rounded-xl gap-4 relative overflow-hidden"
					>
						<div class="flex items-center gap-3">
							<div
								class="w-9 h-9 rounded-full bg-slate-100 border border-slate-200 text-slate-400 flex items-center justify-center shrink-0"
							>
								{#if badge.icon === 'heart'}
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
											d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
										/>
									</svg>
								{:else if badge.icon === 'award'}
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
											d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25z"
										/>
									</svg>
								{/if}
							</div>

							<div>
								<h3 class="text-xs font-bold text-slate-600 leading-tight">{badge.name}</h3>
								<p class="text-[10px] text-slate-400 mt-0.5">{badge.description}</p>
							</div>
						</div>

						<div class="shrink-0 text-slate-400">
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
									d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25z"
								/>
							</svg>
						</div>
					</div>
				{/each}
			</div>
		</section>
	</div>
</div>
