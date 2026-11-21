import sys
import json
from ortools.sat.python import cp_model


def solve_timetable(data):
    model = cp_model.CpModel()

    days_count = data.get("daysCount", 5)
    slots_per_day = data.get("slotsPerDay", 6)
    rooms = data.get("rooms", [])
    blocks = data.get("blocks", [])
    preferences = data.get("preferences", [])

    days = range(1, days_count + 1)
    slots = range(1, slots_per_day + 1)

    if not rooms or not blocks:
        return {"status": "EMPTY_INPUT", "assignments": []}

    # Group room IDs by their roomTypeId for quick lookup during scheduling.
    rooms_by_type = {}
    for r in rooms:
        rooms_by_type.setdefault(r["roomTypeId"], []).append(r["id"])

    # 1. متغيرات القرار الثنائية: x[(b_idx, s_num, room_id, day, slot)]
    x = {}
    for b_idx, blk in enumerate(blocks):
        valid_rooms = rooms_by_type.get(blk["roomTypeId"], [])
        for s_num in range(blk["weeklySessions"]):
            for r_id in valid_rooms:
                for d in days:
                    for s in slots:
                        x[(b_idx, s_num, r_id, d, s)] = model.NewBoolVar(
                            f"b{b_idx}_n{s_num}_r{r_id}_d{d}_s{s}"
                        )

    # 2. القيد 1: جدولة كل حصة مطلوبة مرة واحدة بالضبط
    for b_idx, blk in enumerate(blocks):
        valid_rooms = rooms_by_type.get(blk["roomTypeId"], [])
        if not valid_rooms:
            return {"status": "INFEASIBLE", "assignments": []}
        for s_num in range(blk["weeklySessions"]):
            model.AddExactlyOne(
                x[(b_idx, s_num, r_id, d, s)]
                for r_id in valid_rooms for d in days for s in slots
            )

    # 3. القيد 2: منع تكرار نفس الحصة في نفس اليوم إذا تكررت أسبوعياً
    for b_idx, blk in enumerate(blocks):
        if blk["weeklySessions"] > 1 and blk["weeklySessions"] <= days_count:
            valid_rooms = rooms_by_type.get(blk["roomTypeId"], [])
            for d in days:
                model.AddAtMostOne(
                    x[(b_idx, s_num, r_id, d, s)]
                    for s_num in range(blk["weeklySessions"])
                    for r_id in valid_rooms
                    for s in slots
                )

    # 4. القيد 3: منع تضارب القاعات
    for r in rooms:
        r_id = r["id"]
        for d in days:
            for s in slots:
                vars_in_room = [
                    x[(b_idx, s_num, r_id, d, s)]
                    for b_idx, blk in enumerate(blocks)
                    for s_num in range(blk["weeklySessions"])
                    if (b_idx, s_num, r_id, d, s) in x
                ]
                if vars_in_room:
                    model.AddAtMostOne(vars_in_room)

    # 5. القيد 4: منع تضارب الأساتذة
    prof_ids = set(blk["professorId"] for blk in blocks)
    for p_id in prof_ids:
        for d in days:
            for s in slots:
                vars_for_prof = [
                    x[(b_idx, s_num, r_id, d, s)]
                    for b_idx, blk in enumerate(blocks)
                    if blk["professorId"] == p_id
                    for s_num in range(blk["weeklySessions"])
                    for r_id in rooms_by_type.get(blk["roomTypeId"], [])
                ]
                if vars_for_prof:
                    model.AddAtMostOne(vars_for_prof)

    # 6. القيد 5: منع تضارب الأفواج الطلابية (بما فيها المحاضرات المشتركة)
    all_groups = set(g for blk in blocks for g in blk.get("groupIds", []))
    for g_id in all_groups:
        for d in days:
            for s in slots:
                vars_for_group = [
                    x[(b_idx, s_num, r_id, d, s)]
                    for b_idx, blk in enumerate(blocks)
                    if g_id in blk.get("groupIds", [])
                    for s_num in range(blk["weeklySessions"])
                    for r_id in rooms_by_type.get(blk["roomTypeId"], [])
                ]
                if vars_for_group:
                    model.AddAtMostOne(vars_for_group)

    # 7. تطبيق تفضيلات وقيود الأساتذة المؤقتة
    forbidden_prof_slots = set()
    preferred_prof_bonus = {}

    for pref in preferences:
        p_id = pref["professorId"]
        d = pref["dayNumber"]
        s = pref["slotNumber"]
        score = pref["score"]

        target_slots = slots if s == 0 else [s]
        for target_s in target_slots:
            if score < 0:
                forbidden_prof_slots.add((p_id, d, target_s))
            elif score > 0:
                preferred_prof_bonus[(p_id, d, target_s)] = score

    for (b_idx, s_num, r_id, d, s), var in x.items():
        p_id = blocks[b_idx]["professorId"]
        if (p_id, d, s) in forbidden_prof_slots:
            model.Add(var == 0)

    # 8. دالة الهدف: تفضيل الحصص الصباحية + مكافأة الساعات المفضلة للأستاذ
    objective_terms = []
    for (b_idx, s_num, r_id, d, s), var in x.items():
        p_id = blocks[b_idx]["professorId"]
        bonus = preferred_prof_bonus.get((p_id, d, s), 0)
        cost = s - bonus
        objective_terms.append(var * cost)

    model.Minimize(sum(objective_terms))

    # 9. تشغيل المحرك الرياضي
    # solver = cp_model.CpSolver()
    # solver.parameters.max_time_in_seconds = 15.0
    # status = solver.Solve(model)

    # 9. تشغيل المحرك الرياضي بأقصى سرعة (معالجة متعددة الأنوية + حد زمني 5 ثوانٍ)
    solver = cp_model.CpSolver()
    solver.parameters.max_time_in_seconds = 5.0
    solver.parameters.num_search_workers = 8  # استخدام جميع أنوية المعالج
    status = solver.Solve(model)

    if status not in (cp_model.OPTIMAL, cp_model.FEASIBLE):
        return {"status": "INFEASIBLE", "assignments": []}

    assignments = []
    for (b_idx, s_num, r_id, d, s), var in x.items():
        if solver.Value(var) == 1:
            blk = blocks[b_idx]
            assignments.append({
                "offeringId": blk["offeringId"],
                "sessionUuid": blk["sessionUuid"],
                "roomId": r_id,
                "dayNumber": d,
                "slotNumber": s
            })

    status_label = "OPTIMAL" if status == cp_model.OPTIMAL else "FEASIBLE"
    return {"status": status_label, "assignments": assignments}


if __name__ == "__main__":
    input_data = json.load(sys.stdin)
    result = solve_timetable(input_data)
    print(json.dumps(result))