const std = @import("std");
const types = @import("types.zig");
const testing = @import("std").testing;

const ErrorSet = error{InvalidInput};

pub fn readInput(allocator: std.mem.Allocator, path: types.string) ![]const u8 {
    const file = try std.fs.cwd().openFile(path, .{});
    defer file.close();

    // AoC inputs are usually small (<30KB), so 1MB is plenty of buffer.
    const max_size = 1_000_000;
    const content = try file.readToEndAlloc(allocator, max_size);

    // remove trailing whitespace
    return std.mem.trimRight(u8, content, "\r\n");
}

pub fn iterLines(input: types.string) std.mem.SplitIterator(u8, .scalar) {
    return std.mem.splitScalar(u8, input, '\n');
}

// helper to check if a string is effectively empty
pub fn isEmpty(line: types.string) bool {
    return std.mem.trim(u8, line, "\r\t ").len == 0;
}

fn Point(comptime T: type) type {
    return struct {
        x: isize,
        y: isize,
        val: T,
    };
}

pub fn Grid(comptime T: type, U: type) type {
    return struct {
        data: []const T,
        width: usize,
        height: usize,
        allocator: std.mem.Allocator,

        const Self = @This();

        pub const Error = error{
            InvalidDimensions,
            OutOfMemory,
            ParseFailure,
        };

        pub fn fromLines(
            allocator: std.mem.Allocator,
            lines: []const []const u8,
            comptime parseCell: fn (T) ErrorSet!U,
        ) !Self {
            if (lines.len == 0) return Error.InvalidDimensions;
            const width = lines[0].len;
            for (lines[1..]) |line| {
                if (line.len != width) return Error.InvalidDimensions;
            }

            const data = try allocator.alloc(T, width * lines.len);
            errdefer allocator.free(data);

            for (lines, 0..) |line, y| {
                for (line, 0..) |_, x| {
                    const cell_src = line[x .. x + 1]; // one char
                    data[y * width + x] = try parseCell(cell_src);
                }
            }

            return Self{
                .data = data,
                .width = width,
                .height = lines.len,
                .allocator = allocator,
            };
        }

        pub fn deinit(self: Self) void {
            self.allocator.free(self.data);
        }

        pub fn at(self: Self, x: usize, y: usize) ?Point(T) {
            if (x >= self.width or y >= self.height or x < 0 or y < 0) return null;
            return Point(T){ .x = @intCast(x), .y = @intCast(y), .val = self.data[y * self.width + x] };
        }

        /// Helper for "get neighbors"
        pub fn neighbors(self: Self, x: usize, y: usize) struct { north: ?Point(T), south: ?Point(T), east: ?Point(T), west: ?Point(T) } {
            const north: ?Point(T) = at(self, x, y - 1);
            const east: ?Point(T) = at(self, x + 1, y);
            const south: ?Point(T) = at(self, x, y + 1);
            const west: ?Point(T) = at(self, x - 1, y);

            return .{
                .north = north,
                .south = south,
                .east = east,
                .west = west,
            };
        }
    };
}

fn parse(input: []const u8) ErrorSet![]const u8 {
    return input;
}

test "Grid.at" {
    const allocator = std.testing.allocator;
    const lines = [_][]const u8{
        "123",
        "456",
        "789",
    };

    const grid = try Grid([]const u8, []const u8).fromLines(
        allocator,
        &lines,
        parse,
    );
    defer grid.deinit();

    const expected: []const u8 = "4";
    const actual: []const u8 = grid.at(0, 1).?.val;
    try testing.expectEqualStrings(expected, actual);
}

test "Grid.neighbors" {
    const allocator = std.testing.allocator;
    const lines = [_][]const u8{
        "123",
        "456",
        "789",
    };
    const grid = try Grid([]const u8, []const u8).fromLines(allocator, &lines, parse);
    defer grid.deinit();

    const expected = struct {
        north: ?Point([]const u8),
        south: ?Point([]const u8),
        east: ?Point([]const u8),
        west: ?Point([]const u8),
    }{
        .north = Point([]const u8){ .x = 1, .y = 0, .val = "2" },
        .south = Point([]const u8){ .x = 1, .y = 2, .val = "8" },
        .east = Point([]const u8){ .x = 2, .y = 1, .val = "6" },
        .west = Point([]const u8){ .x = 0, .y = 1, .val = "4" },
    };
    const actual = grid.neighbors(1, 1);

    try testing.expectEqualDeep(expected.north, actual.north);
    try testing.expectEqualDeep(expected.south, actual.south);
    try testing.expectEqualDeep(expected.east, actual.east);
    try testing.expectEqualDeep(expected.west, actual.west);
}
