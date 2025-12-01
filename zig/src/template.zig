const std = @import("std");
const utils = @import("utils.zig");

pub fn solve(allocator: std.mem.Allocator) !void {
    // 1. Read Input
    // Note: We read from the "inputs" folder relative to where we run the command
    const input = try utils.readInput(allocator, "inputs/day01.txt");
    defer allocator.free(input);

    // 2. Parse (Common for Part 1 & 2)
    var lines = utils.iterLines(input);

    // 3. Solve
    var part1_result: u64 = 0;
    var part2_result: u64 = 0;

    while (lines.next()) |line| {
        if (line.len == 0) continue;
        // Do logic here...
        part1_result += 1;
        part2_result += 1;
    }

    // 4. Print Answers
    std.debug.print("Day 01 Part 1: {d}\n", .{part1_result});
    std.debug.print("Day 01 Part 2: {d}\n", .{part2_result});
}
