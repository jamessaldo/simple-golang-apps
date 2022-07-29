//
//  Model.swift
//  Tangled
//
//  Created by zy on 28/07/22.
//

import Foundation

struct Post: Codable {
    var id: Int?
    var title: String?
    var description: String?
    var creator: String?
    var created_at: String?
    var updated_at: String?
    var comments: [Comment]?
}

struct CreatePost: Codable {
    var title: String
    var description: String
}

struct CreateComment: Codable {
    var content: String
    var post_id: Int
}


struct Comment: Codable {
    var id: Int?
    var post_id: Int?
    var content: String?
    var creator: String?
    var created_at: String?
    var updated_at: String?
}
