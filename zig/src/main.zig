const std = @import("std");
const utils = @import("utils.zig");

// Import your days here
// const day01 = @import("day01.zig");
// const day02 = @import("day02.zig");
// const day03 = @import("day03.zig");

pub fn main() !void {
    // 1. Setup Allocator (GPA is best for debugging leaks)
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    _ = try utils.readInput(allocator, "abc123");

    // 2. Parse Args to get the day number
    const args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len < 2) {
        std.debug.print("Usage: zig build run -- <day_number>\n", .{});
        return;
    }

    const day_num = try std.fmt.parseInt(u8, args[1], 10);

    // 3. Dispatch to the correct day
    // We time the execution because speed is fun!
    const start = std.time.nanoTimestamp();

    switch (day_num) {
        // 1 => try day01.solve(allocator),
        // 2 => try day02.solve(allocator),
        // 3 => try day03.solve(allocator),
        else => std.debug.print("Day {d} not implemented yet!\n", .{day_num}),
    }

    const end = std.time.nanoTimestamp();
    const elapsed = @as(f64, @floatFromInt(end - start)) / 1_000_000.0;
    std.debug.print("\nTime: {d:.3}ms\n", .{elapsed});
}
