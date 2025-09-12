package algos

import "context"

func cycleSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()

	for cycleStart := 0; cycleStart < n-1; cycleStart++ {
		item := arr.Get(ctx, cycleStart)
		pos := cycleStart

		for i := cycleStart + 1; i < n; i++ {
			if arr.Get(ctx, i) < item {
				pos++
			}
		}

		if pos == cycleStart {
			continue
		}

		for arr.Get(ctx, pos) == item {
			pos++
		}

		if pos != cycleStart {
			temp := arr.Get(ctx, pos)
			arr.Set(ctx, pos, item)
			item = temp
		}

		for pos != cycleStart {
			pos = cycleStart

			for i := cycleStart + 1; i < n; i++ {
				if arr.Get(ctx, i) < item {
					pos++
				}
			}

			for arr.Get(ctx, pos) == item {
				pos++
			}

			if arr.Get(ctx, pos) != item {
				temp := arr.Get(ctx, pos)
				arr.Set(ctx, pos, item)
				item = temp
			}
		}
	}
}
